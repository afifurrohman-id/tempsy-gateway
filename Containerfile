FROM rust:1.76-alpine AS builder
WORKDIR /src

# Need install ca-certificates for tls compatibility for go library 
# hadolint ignore=DL3018
RUN apk add --no-cache \
  musl-dev
# ca-certificates && \
# update-ca-certificates

COPY . .

ENV CGO_ENABLED=0

# Reduce binary size by removing debug information
RUN cargo fix --allow-dirty --allow-staged && \
  cargo build -r

FROM scratch
LABEL org.opencontainers.image.authors="afif"
LABEL org.opencontainers.image.licenses="MIT"
WORKDIR /app

COPY --from=builder /src/target/release/gateway .
# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD [ "./gateway" ]
