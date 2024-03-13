FROM lukemathwalker/cargo-chef:latest-rust-1.76-slim AS chef
WORKDIR /app

FROM chef AS planner
COPY . .
RUN cargo chef prepare --recipe-path recipe.json

FROM chef AS builder 

# hadolint ignore=DL3008
RUN apt-get update && \
  apt-get install \ 
  --no-install-recommends \
  perl make -y

COPY --from=planner /app/recipe.json recipe.json
RUN cargo chef cook --release \ 
  --recipe-path recipe.json
COPY . .
RUN cargo build -r

FROM debian:bookworm-slim
LABEL org.opencontainers.image.authors="afif"
LABEL org.opencontainers.image.licenses="MIT"
WORKDIR /app

COPY --from=builder /app/target/release/gateway .
RUN ls /app

CMD [ "./gateway" ]
