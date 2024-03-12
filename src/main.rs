#![forbid(unsafe_code)]
// #![warn(clippy::expect_used)]
#![warn(clippy::unwrap_used)]
#![cfg_attr(
    not(debug_assertions),
    forbid(unreachable_code),
    forbid(unused_variables)
)]

use std::env;
use tokio::net as tokio_net;
use tracing_subscriber::{fmt, layer::SubscriberExt, util::SubscriberInitExt};

#[tokio::main]
async fn main() {
    #[cfg(debug_assertions)]
    dotenvy::dotenv().unwrap_or_else(|err| panic!("Failed to load `.env` file: {err}"));

    let host = env::var("HOST").unwrap_or_else(|_| {
        let host = "0.0.0.0".to_string();
        eprintln!("`HOST` env is not set, using `{host}` for fallback");
        host
    });
    let port = env::var("PORT").expect("`PORT` env must be set");

    tracing_subscriber::registry()
        .with(fmt::layer().pretty())
        .init();

    let lis = tokio_net::TcpListener::bind(&format!("{host}:{port}"))
        .await
        .unwrap_or_else(|err| panic!("Failed to bind {host}:{port}, error: {err}"));

    tracing::info!(
        "App listen on: {}",
        lis.local_addr().expect("Local address must be available")
    );

    axum::serve(lis, tempsy_gateway::router())
        .await
        .unwrap_or_else(|err| panic!("Failed to run app: {err}"))
}
