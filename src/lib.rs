#![forbid(unsafe_code)]
// #![warn(clippy::expect_used)]
#![warn(clippy::unwrap_used)]
#![cfg_attr(
    not(debug_assertions),
    forbid(unreachable_code),
    forbid(unused_variables)
)]

use axum::http::header;
use axum::response::IntoResponse;
use axum::{body, extract, http, routing};
use tower_http::cors;

use hyper_util::client::legacy;
use hyper_util::client::legacy::connect;
use hyper_util::rt;
use std::{any, env};
use tower_http::{catch_panic, trace};

type Client = legacy::Client<hyper_tls::HttpsConnector<connect::HttpConnector>, body::Body>;
const HEADER_FILE_IS_PUBLIC: &str = "file-is-public";
const HEADER_FILE_PRIVATE_URL_EXPIRES: &str = "file-private-url-expires";
const HEADER_FILE_AUTO_DELETE_AT: &str = "file-auto-delete-at";
const HEADER_FILE_NAME: &str = "file-name";

async fn gateway_handler(
    extract::State(client): extract::State<Client>,
    mut req: extract::Request,
) -> Result<impl IntoResponse, String> {
    let path = req.uri().path();
    let path_query = req
        .uri()
        .path_and_query()
        .map(|pq| pq.as_str())
        .unwrap_or(path);

    let mut var = "CLIENT_URL";

    if regex::Regex::new(r"^/files/[a-zA-Z0-9_-]+/public/[a-zA-Z0-9_-]+\.+[a-zA-Z0-9_-]+$")
        .unwrap_or_else(|err| panic!("Cannot create regex: {err}"))
        .is_match(path)
    {
        var = "SERVER_URL";
    }

    if let Some(acc_type) = req
        .headers()
        .get(header::ACCEPT)
        .and_then(|ct| ct.to_str().ok())
    {
        if acc_type.contains("application/json") {
            var = "SERVER_URL";
        }
    }

    let uri = format!(
        "{}{}",
        env::var(var).unwrap_or_else(|err| panic!("could not get env {var}: {err}")),
        path_query
    );

    req.headers_mut().insert(header::HOST, {
        let host = uri
            .trim_end_matches('/')
            .split("://")
            .nth(1)
            .unwrap_or_default();
        host.parse()
            .unwrap_or_else(|err| panic!("HOST: `{host}`, is invalid for header value: {err}"))
    });

    *req.uri_mut() = hyper::Uri::try_from(uri).unwrap_or_default();
    Ok(client
        .request(req)
        .await
        .unwrap_or_else(|err| panic!("Error make request to server: {err}"))
        .into_response())
}

pub fn router() -> axum::Router {
    axum::Router::new()
        .route("/", routing::any(gateway_handler))
        .fallback(routing::any(gateway_handler))
        .layer(trace::TraceLayer::new_for_http())
        .layer(catch_panic::CatchPanicLayer::custom(handle_panic))
        .layer(
            cors::CorsLayer::new()
                .allow_origin(cors::Any)
                .allow_methods([
                    http::Method::GET,
                    http::Method::POST,
                    http::Method::PUT,
                    http::Method::DELETE,
                    http::Method::OPTIONS,
                ])
                .allow_headers([
                    header::CONTENT_TYPE,
                    header::AUTHORIZATION,
                    header::HeaderName::from_static(HEADER_FILE_AUTO_DELETE_AT),
                    header::HeaderName::from_static(HEADER_FILE_PRIVATE_URL_EXPIRES),
                    header::HeaderName::from_static(HEADER_FILE_IS_PUBLIC),
                    header::HeaderName::from_static(HEADER_FILE_NAME),
                ]),
        )
        .with_state(client())
}

fn client() -> Client {
    legacy::Client::<(), ()>::builder(rt::TokioExecutor::new())
        .build(hyper_tls::HttpsConnector::new())
}

fn handle_panic(err: Box<dyn any::Any + Send + 'static>) -> hyper::Response<String> {
    if let Some(err) = err.downcast_ref::<String>() {
        tracing::error!("{err}")
    }
    hyper::Response::builder()
        .status(hyper::StatusCode::INTERNAL_SERVER_ERROR)
        .body("Internal server error".to_string())
        .unwrap_or_default()
}
