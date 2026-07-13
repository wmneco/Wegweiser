## Why

Wegweiser v0.1.0 needs two endpoints — shorten a URL and redirect a short code — but before that behaviour can be built we need HTTP-reachable service binaries with consistent, testable plumbing. This change delivers that groundwork only: binaries you can `curl` and get a well-formed response from, so the later endpoint and storage issues stay focused on behaviour rather than transport.

## What Changes

- Introduce **two independently deployable service binaries** — `redirectd` and `shortend` — each reachable over HTTP with graceful startup/shutdown.
- Structure each service as a **hexagon (ports & adapters)**: an HTTP driving adapter calls an inbound port backed by a core type whose method is a placeholder for now. Only the driving side is built; the outbound store port and its adapter are deferred to the storage issue.
- Register **placeholder routes** so both endpoints are reachable but not yet functional: `GET /{code}` (redirect) and `POST /api/shorten` (shorten).
- Add a **shared runtime** used by both binaries for mechanics only (nothing that knows about URLs): server lifecycle, baseline middleware (request ID, panic recovery, request logging), a liveness health endpoint, a consistent error-response shape, structured logging, and environment-based configuration with sane overridable defaults.
- Establish **build/run/test/lint entrypoints** and a smoke test against the health endpoint.

Out of scope (separate v0.1.0 issues): actual shorten/redirect logic, the outbound store ports and adapters, any persistent storage backend, auth, rate limiting, and custom aliases.

## Capabilities

### New Capabilities
- `service-runtime`: Shared HTTP runtime mechanics reused by every service binary — graceful lifecycle, request-ID / panic-recovery / request-logging middleware, a liveness health endpoint, a consistent error-response shape, structured logging, and environment-based configuration. Contains no URL-shortener domain knowledge.
- `redirect-service`: The `redirectd` binary — boots on its configured address, composes the shared runtime, and exposes a placeholder `GET /{code}` route wired through an HTTP adapter → inbound port → core stub.
- `shorten-service`: The `shortend` binary — boots on its configured address, composes the shared runtime, and exposes a placeholder `POST /api/shorten` route wired through an HTTP adapter → inbound port → core stub.

### Modified Capabilities
<!-- None — this is the first code in the repository; no existing specs change. -->

## Impact

- **New packages**: `cmd/redirectd`, `cmd/shortend`; `internal/platform/*` (shared runtime); `internal/redirect/*` and `internal/shorten/*` (per-service hexagons, driving side only).
- **Dependencies**: production code stays standard-library only (`net/http` method/pattern routing, `log/slog`, `os`/env parsing); `github.com/stretchr/testify` is added as a **test-only** dependency for assertion ergonomics. Discussed per the "discuss deps before adding" rule; rationale recorded in `design.md` (D8).
- **Deferred coupling**: splitting into two processes now removes the in-memory-map option from the future storage issue; that issue must open with a shared/networked store. Accepted deliberately.
- **CI/docs**: existing `go.yml` (gofmt, goimports, vet, `go test -race` + coverage) already covers the new code; `CHANGELOG.md` `[Unreleased]` gets an entry.
