# SSO

Authentication service for [K-Tify](https://github.com/DavidJackso/K-Tify) — a microservice music-streaming platform. Handles user registration and login over **gRPC** and issues **JWT** access tokens.

**Stack:** Go 1.24 · gRPC · PostgreSQL 16 · bcrypt · golang-jwt · cleanenv · slog · Docker Compose · golang-migrate · Taskfile

## gRPC API

Proto contracts: [`K-Tify/protos/protos/sso/sso.proto`](https://github.com/DavidJackso/K-Tify/tree/main/protos/protos/sso).

| Method | Request | Response | Description |
|---|---|---|---|
| `Auth.Register` | `email`, `password` | `user_id` | Hashes the password with bcrypt and creates the user |
| `Auth.Login` | `email`, `password` | `token` | Verifies credentials and returns a signed JWT |

Invalid input returns `codes.InvalidArgument`.

**JWT:** HS256, signed with the client app's secret (`apps` table). Claims: `uid`, `email`, `exp`. TTL is set by `token_ttl`.

## Architecture

```
cmd/main.go                  entrypoint: config → logger → app, graceful shutdown on SIGINT
internal/
├── app/                     wiring: storage → service → gRPC server
│   └── grpc/                gRPC server lifecycle
├── grpc/auth/               transport layer: request validation, status codes
├── services/auth/           business logic: register, login, bcrypt
│   └── lib/jwt/             token issuing
├── storage/pg/              PostgreSQL repository
├── domain/models/           User, App
└── config/                  YAML config (cleanenv)
migrations/                  SQL migrations (users, apps)
```

Layers depend on interfaces (`UserSaver`, `UserProvider`, `AppProvider`), so storage can be swapped or mocked in tests.

## Quick start

Requires Docker and [Task](https://taskfile.dev).

```bash
task rebuild          # reset containers, start Postgres (localhost:5433), run migrations

# register a client app whose secret signs tokens
docker exec -it auth-db psql -U postgres -d authdb \
  -c "INSERT INTO apps (id, name, secret) VALUES (1, 'k-tify', 'change-me');"

CONFIG_PATH=./config/local.yaml go run ./cmd
```

The service listens on `:44045`.

Try it with [grpcurl](https://github.com/fullstorydev/grpcurl) (server reflection is off, so pass the proto file):

```bash
grpcurl -plaintext -import-path ../K-Tify/protos/protos/sso -proto sso.proto -d '{"email":"user@example.com","password":"secret"}' localhost:44045 auth.Auth/Register
grpcurl -plaintext -import-path ../K-Tify/protos/protos/sso -proto sso.proto -d '{"email":"user@example.com","password":"secret"}' localhost:44045 auth.Auth/Login
```

## Configuration

`config/local.yaml`, path passed via `CONFIG_PATH`:

| Key | Default | Description |
|---|---|---|
| `env` | `local` | `local` — text logs (debug), `dev` — JSON logs (info) |
| `token_ttl` | `1h` | JWT lifetime |
| `grpc.port` | `4041` | gRPC port (`44045` in local.yaml) |
| `grpc.timeout` | `5s` | Request timeout |
| `db.host` / `db.port` | `localhost` / `5432` | PostgreSQL address |
| `db.user` / `db.password` / `db.dbname` | `postgres` | Credentials |

## Tasks

| Command | Description |
|---|---|
| `task db-up` | Start PostgreSQL only |
| `task migrate` | Apply migrations |
| `task rebuild` | Reset, start DB, migrate |
| `task logs` / `task ps` | Logs / container status |
| `task reset` | Remove containers and volumes |

## Roadmap

- App ID in `Login` request instead of the default app
- Map domain errors to gRPC codes (`NotFound`, `AlreadyExists`, `Unauthenticated`)
- Unit tests for service and JWT layers
- `IsAdmin` / token validation RPC
