## Context

Wegweiser is a self-hosted URL shortener with an explicit interest in exploring cloud-native engineering, including eventual scaling of the redirect path. The read path (redirect) and write path (shorten) are strongly asymmetric: redirects are the vast majority of traffic, latency-critical, cacheable, and must stay highly available, while shortens are low-volume writes that tolerate more latency. This asymmetry is the standard justification for eventually running them as separate services.

This change is the transport-only groundwork under the v0.1.0 milestone (issue #1). It delivers HTTP-reachable binaries with consistent plumbing and placeholder handlers. It does **not** deliver the shorten/redirect logic or any storage — those are separate v0.1.0 issues.

Constraints:
- Standard library only; new dependencies require prior discussion (per AGENTS.md).
- Go 1.26 — `net/http.ServeMux` supports method + path-pattern routing, covering this app's needs.
- CI already enforces gofmt, goimports, `go vet`, and `go test -race` with coverage, so the design must be testable without a live network (dependencies injected, no globals).

## Goals / Non-Goals

**Goals:**
- Two independently deployable binaries (`redirectd`, `shortend`) reachable over HTTP with graceful lifecycle.
- Each service structured as a hexagon (ports & adapters) with its driving side wired now.
- Shared runtime for mechanics only, so both binaries stand up an identical operational surface without copy-paste.
- Placeholder routes reachable and returning a well-formed response.
- Build/run/test/lint entrypoints and a health smoke test.

**Non-Goals:**
- Shorten/redirect behaviour (later issues).
- Outbound store ports, store adapters, and any persistent storage backend (later issues).
- Auth, rate limiting, custom aliases, caching, and cross-service coordination.

## Decisions

### D1: Two service binaries from day one

Ship `cmd/redirectd` and `cmd/shortend` as separate binaries now rather than a single binary split later.

- **Why**: aligns with the cloud-native goal of scaling the read path independently; and at the scaffolding stage there is no shared state to coordinate, so the split is cheap here.
- **Consequence (accepted)**: separating into two processes removes the "simple in-memory map" option from the future storage issue — that issue must open with a shared/networked store. This is a deliberate trade, spent now to buy the split.
- **Bonus**: the split dissolves a routing wart. In a monolith, `GET /{code}` at the root competes with `/api/shorten` and `/healthz` for the namespace; split apart, the redirect service owns its entire root and its only other route is health.
- **Alternative considered**: modular monolith with a pre-cut seam, splitting at a later scaling issue. Rejected because the project explicitly wants the split and the milestone treats storage as a separate concern anyway.

### D2: Hexagonal architecture, one hexagon per service

Each service is its own bounded context / hexagon. Dependencies point inward: adapters depend on the core; the core depends on nothing external.

```
   ┌──────── redirect hexagon ────────┐   ┌──────── shorten hexagon ─────────┐
   │ HTTP GET /{code} ─▶ [Resolver]    │   │ HTTP POST /api/shorten ─▶         │
   │  (driving adapter)  (core, stub)  │   │  (driving adapter)  [Shortener]   │
   └───────────────┬───────────────────┘   └───────────────┬──────────────────┘
                   └───────── both built from ──────────────┘
                                    ▼
                 shared runtime: server · middleware · health · respond · config
```

- **Why hexagonal**: the entire cloud-native roadmap (in-mem → Redis/Postgres, HTTP → possibly gRPC) is a series of adapter swaps around a stable core. Ports & adapters makes those swaps local. It also maps cleanly onto idiomatic Go, where consumer-defined interfaces are already the house style.

### D3: Build the driving side only now (half-hexagon), with a real core stub

At this issue there is no domain logic and no storage, so only the driving spine is built: HTTP adapter → inbound port → core type whose method returns a placeholder. The outbound store port and its driven adapter are deferred to the storage issue, where they will have logic to serve.

- **Chosen (option A)**: the core exists now as a named type implementing the inbound port, with a placeholder method body. The follow-up logic issue only fills the method; wiring is already in place.
- **Alternative (option B)**: skip the core entirely and have the HTTP adapter return the stub directly. Rejected — it defers wiring into the logic issue, turning scaffolding into something that must be rearranged later. Option A makes the scaffolding *be* the architecture at the cost of a few lines.
- **Guardrail**: do not fabricate outbound ports for behaviour that does not exist yet; the hexagon grows its driven side when storage arrives.

### D4: Shared runtime is mechanics only — "share mechanics, never share core"

A shared `internal/platform` package holds everything that knows nothing about URLs; anything reasoning about codes or targets stays inside a specific hexagon.

| Shared (`internal/platform`) | Never shared (inside each hexagon) |
| --- | --- |
| server lifecycle / graceful shutdown | domain core (`Resolver`, `Shortener`) |
| middleware: request ID, panic recovery, logging | inbound ports (the app API) |
| liveness health handler | outbound ports (`Lookup`, `Save`) — later |
| JSON error-response writer | which error → which status code |
| config parsing mechanism | what each config value means |

- **Health nuance**: liveness ("process is up") is a pure mechanic and shared. A future *readiness* check that probes a store touches a driven port and would become per-service; out of scope here.

### D5: Segregated, per-service outbound ports (recorded now, built later)

When storage arrives, each core defines the outbound port it needs — `Lookup(ctx, code) (url, error)` for redirect, `Save(ctx, code, url) error` for shorten — rather than a shared `Store` interface. Interface segregation falls out for free and matches Go's "define the interface where you consume it" idiom. Recorded here so the later issue does not regress into a shared interface.

### D6: Standard library for production code

Production code uses the standard library — this is a deliberate fit for this change's small surface, not merely the "discuss deps first" constraint. See D8 for the one test-only dependency and the general dependency policy.

- **Routing**: `net/http.ServeMux` with method + pattern routes (Go 1.22+). No third-party router — the app has ~2–3 routes per service and hexagonal keeps routing thin, so a router would mostly wrap stdlib for near-zero gain.
- **Logging**: `log/slog`, structured, configurable level.
- **Config**: environment variables via `os` with a small parse function and documented defaults (listen address, log level, shutdown timeout). 12-factor-friendly; each binary reads its own env in its own container.
- **Composition**: the shared server is a plain runner the composition root (`main`) calls after building adapters and wiring the core — favouring the boring function over a configurable object.

### D8: Dependency policy — stdlib for production, testify for tests

**Decision**: production code stays stdlib-only for this change; add `github.com/stretchr/testify` (`assert`/`require`) as a **test-only** dependency.

**Why testify**: it is the de-facto standard for Go test assertions, and the project mandates tests with coverage and `-race`, so test readability compounds across the suite. It stays out of production binaries (imported only from `_test.go` files).

**Rejected now, recorded for later** (each cheap to adopt under hexagonal when it earns its place):

- **Third-party router** (chi/gin/echo): ServeMux covers current needs; revisit only if routing grows complex.
- **Config library** (caarlos0/env, envconfig): three env vars do not justify it; reconsider when configuration surface grows.
- **UUID library** (google/uuid): request IDs are ~5 lines over `crypto/rand`.
- **Faster logger** (zerolog/zap): `slog` is right for v0.1.0; swapping the logger is a clean adapter change to make in a dedicated scaling issue when the redirect path is actually hot.

**Policy going forward**: new dependencies are discussed and their trade-offs recorded before adoption (per AGENTS.md); prefer stdlib until a dependency's value is concrete.

### D7: Placeholder responses

Placeholder routes return HTTP 501 Not Implemented using the shared error-response shape, signalling "reachable but not yet implemented." This keeps reachability honestly testable and exercises the error-writer and middleware end to end.

## Risks / Trade-offs

- **Premature split cost** → Accepted in D1: forecloses the in-memory-map storage option. Mitigation: recorded explicitly so the storage issue plans for a shared store from the start.
- **Ceremony risk** — hexagonal structure for placeholder handlers can look like over-engineering → Mitigation: D3 builds only the driving side and forbids fabricating outbound ports; each `main` stays ~15 lines.
- **Drift between the two services' plumbing** → Mitigation: D4 centralises all mechanics in `internal/platform`; services differ only in their route and core.
- **ServeMux pattern edge cases** (e.g. `/{code}` matching `/healthz`) → Largely mitigated by D1 (redirect owns its own root); covered by the reachability tests per D7.
- **Config divergence across binaries** → Same parsing code, per-binary env values; documented defaults keep local runs frictionless.

## Open Questions

- Exact env var names and default values (e.g. `WEGWEISER_ADDR` / `PORT`, `LOG_LEVEL`, `SHUTDOWN_TIMEOUT`) — to settle during implementation.
- Whether to honour an inbound request-ID/correlation header or always generate — leaning toward "use inbound if present, else generate."
- Error-envelope field names (`code`, `message`, optionally `request_id`) — to finalise when writing the shared error writer.
