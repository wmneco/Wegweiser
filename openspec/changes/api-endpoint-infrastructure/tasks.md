## 1. Shared runtime (`internal/platform`)

- [x] 1.1 Add config: parse listen address, log level, and shutdown timeout from environment with documented defaults
- [x] 1.2 Add structured logging setup (`log/slog`) honoring the configured level
- [x] 1.3 Add the error-response writer producing the consistent JSON error shape (code, message) with the matching status code
- [x] 1.4 Add middleware: request ID (use inbound correlation header if present, else generate), panic recovery (returns 500 via the error writer), and request logging (method, path, status, duration, request ID)
- [x] 1.5 Add the liveness health handler responding 200 for `GET /healthz`
- [x] 1.6 Add the server runner: start on the configured address, apply middleware, and shut down gracefully on SIGINT/SIGTERM within the shutdown timeout

## 2. Redirect service (`internal/redirect`, `cmd/redirectd`)

- [x] 2.1 Define the inbound port and a core type implementing it with a placeholder resolve method
- [x] 2.2 Add the HTTP driving adapter for `GET /{code}` that invokes the inbound port and returns 501 via the shared error writer
- [x] 2.3 Add `cmd/redirectd/main.go` composition root: load config, build the core + adapter, register `/{code}` and the health route, run via the shared server runner

## 3. Shorten service (`internal/shorten`, `cmd/shortend`)

- [x] 3.1 Define the inbound port and a core type implementing it with a placeholder shorten method
- [x] 3.2 Add the HTTP driving adapter for `POST /api/shorten` that invokes the inbound port and returns 501 via the shared error writer
- [x] 3.3 Add `cmd/shortend/main.go` composition root: load config, build the core + adapter, register `/api/shorten` and the health route, run via the shared server runner

## 4. Tests

- [x] 4.0 Add `github.com/stretchr/testify` as a test-only dependency (`go get`, tidy) and use `require`/`assert` in the tests below
- [x] 4.1 Health smoke test: `GET /healthz` returns 200 for the runtime handler
- [x] 4.2 Redirect placeholder test: `GET /{code}` is routed and returns 501 with the standard error shape
- [x] 4.3 Shorten placeholder test: `POST /api/shorten` is routed and returns 501 with the standard error shape
- [x] 4.4 Middleware tests: request ID present on responses; panic in a handler is recovered as a 500 in the error shape
- [x] 4.5 Config test: defaults applied when env unset; env overrides defaults

## 5. Entrypoints & docs

- [x] 5.1 Add build/run/test/lint entrypoints (e.g. Makefile targets) mirroring the CI checks (gofmt, goimports, vet, `go test -race`)
- [x] 5.2 Verify `go build ./...`, `go vet ./...`, and `go test -race ./...` pass locally
- [x] 5.3 Add a `CHANGELOG.md` `[Unreleased]` entry for the API endpoint infrastructure
