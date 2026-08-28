# Appdroid server

APK hosting and distribution server written in Go, using Clean Architecture and
DDD. Replaces the previous Bun + Elysia implementation with identical public
API behavior.

## Stack

- [chi](https://github.com/go-chi/chi) — HTTP router
- [sarulabs/di](https://github.com/sarulabs/di) — dependency injection
- [avast/apkparser](https://github.com/avast/apkparser) — pure-Go APK parsing
  (no `aapt` binary required)
- [go.etcd.io/bbolt](https://github.com/etcd-io/bbolt) — embedded key/value store
  (no CGO, `CGO_ENABLED=0` compatible)
- [html/template](https://pkg.go.dev/html/template) — server-rendered pages

## Layout

```
cmd/server/                 composition root: container wiring + HTTP server
internal/container/         sarulabs/di definitions (single builder)
internal/config/            PORT, DATA_DIR env config
internal/domain/app/        App entity, value types, AppRepository port
internal/application/app/   use cases (UploadApk, List, Get, Delete, Serve)
internal/infrastructure/
  apkparser/                APK metadata extraction adapter
  apkstorage/               temp + persistent file storage
  repository/               bbolt adapter for AppRepository
  http/                     chi handlers, routes, SSR page templates
  http/static/              admin frontend (vanilla JS, embedded into the binary)
data/                       runtime data: appdroid.db + uploaded APKs
```

## Development

Requires Go 1.26+.

```bash
just            # gofmt + go vet + tests + build (default recipe)
just test       # go test ./...
just run        # go run ./cmd/server (PORT=3000, DATA_DIR=./data)
just build      # build ./bin/appdroid
just build-static # static binary without CGO
```

## Configuration

| Env var     | Default  | Description                             |
|-------------|----------|-----------------------------------------|
| `PORT`      | `3000`   | HTTP listen port                        |
| `DATA_DIR`  | `./data` | Directory for bbolt db and APK files   |
| `API_KEY`   | `""`     | Shared secret for `X-API-Key` header. When empty the server fails closed and rejects every mutating request (POST/PATCH/DELETE) |

> **Auth in dev.** When `SERVER_MODE` is not `release` (e.g. local `just run`), the
> server also loads `./.env` to seed environment variables. Real env
> variables always win over values in the `.env` file, so `API_KEY` set via
> `docker run -e` or compose `env_file` is never overridden. Add `API_KEY` to
> `.env.example` (copy to `.env`) for local development.

## Caching

- HTML entry documents (`/`, SSR pages) always send `Cache-Control: no-cache`
  and are revalidated on every request.
- JS/CSS static assets send `Cache-Control: public, max-age=86400` in release
  mode (`SERVER_MODE=release`, the Docker image default) and are cached for a
  day.
- In dev mode (`SERVER_MODE` unset, e.g. local `just run`) all static assets
  send `no-cache`, so a local rebuild is picked up immediately.

## API

| Route                       | Behavior                                   |
|-----------------------------|--------------------------------------------|
| `GET /`                     | static admin frontend                      |
| `GET /api/ping`             | `pong` (health check)                      |
| `GET /api/list`             | JSON list of apps                          |
| `POST /api/upload`          | multipart `file` (cap 256 MiB), needs `X-API-Key` |
| `DELETE /api/:id`           | delete app + stored APK, needs `X-API-Key` |
| `GET /file/:file`           | serve stored APK                           |
| `GET /apk/:file/:hash`      | serve stored APK                           |
| `GET /apps`                 | SSR page with Obtainium deep links         |
| `GET /app/:id`              | SSR app page with download link            |

Mutating routes (`POST /api/upload`, `DELETE /api/:id`) require a valid
`X-API-Key` header matching `API_KEY`. Without it the server fails closed and
returns `401`. All `GET` routes (downloads, public pages) stay public.

## Docker

Multi-stage build (`golang:1.26-alpine` → `scratch`) produces a single static
binary with the admin frontend embedded. `data/` is a runtime volume — mount it
to keep the bbolt db and APKs across redeploys.

```bash
docker build -t appdroid .
docker run -p 3000:3000 -v appdroid-data:/app/data appdroid
```
