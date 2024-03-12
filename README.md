# [Tempsy Gateway](https://tempsy.afifurrohman.my.id)

> A gateway for [tempsy project](https://github.com/afifurrohman-id/tempsy.git)

## Usage

### Requirements

- [x] WSL2 (Windows Subsystem for Linux)
  > Only need if you're using Windows
- [x] Git (version >= 2.39.x)
- [x] Rust toolchains (version >= 1.76.x)
- [ ] Docker (version >= 24.0.x)
  > Optional, Only need if you want to build image

### Installation

- Clone this repository

```sh
git clone https://github.com/afifurrohman-id/tempsy-gateway.git
```

- Go to the project directory

```sh
cd tempsy-gateway
```

- Insert variables into `.env` file

```sh
cat <<EOENV > configs/.env
APP_ENV=development
PORT=8080
SERVER_URL=https://example.com
CLIENT_URL=https://www.example.com
EOENV
```

### Run

- Run Server

```sh
cargo run
```

- Build

```sh
cargo build -r
```

- Build Image

```sh
docker build -t tempsy-gateway -f Containerfile .
```
