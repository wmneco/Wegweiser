## ADDED Requirements

### Requirement: Redirect service binary is HTTP-reachable

The `redirectd` binary SHALL build, run locally, and serve HTTP on its configured address, composing the shared service runtime for lifecycle, middleware, health, and configuration.

#### Scenario: Binary boots and serves

- **WHEN** `redirectd` is started with a valid configuration
- **THEN** it binds its configured address and begins serving HTTP requests

#### Scenario: Health endpoint reachable

- **WHEN** a client sends `GET /healthz` to a running `redirectd`
- **THEN** it responds with HTTP 200 indicating the service is alive

### Requirement: Placeholder redirect route

The `redirectd` binary SHALL register a placeholder route `GET /{code}` that is reachable but not yet functional. The placeholder SHALL return a well-formed response using the shared error-response shape indicating the behaviour is not yet implemented.

#### Scenario: Placeholder route reachable

- **WHEN** a client sends `GET /{code}` for any code
- **THEN** the route is matched and returns HTTP 501 with the standard error-response shape indicating redirect is not yet implemented

### Requirement: Hexagonal driving-side wiring

The redirect service SHALL be structured as a hexagon whose driving side is wired now: the HTTP adapter SHALL depend on an inbound port, backed by a core type whose resolving method is a placeholder. The outbound store port and its adapter are intentionally deferred to a later change.

#### Scenario: Request flows through the driving side

- **WHEN** the placeholder route handles a request
- **THEN** the HTTP adapter invokes the inbound port implemented by the core type, rather than embedding behaviour directly in the handler
