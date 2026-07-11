## ADDED Requirements

### Requirement: Graceful server lifecycle

The service runtime SHALL provide a reusable way to start an HTTP server on a configured address and shut it down gracefully. On receiving an interrupt or termination signal, the runtime SHALL stop accepting new connections and allow in-flight requests to complete within a configured timeout before exiting.

#### Scenario: Server starts and serves

- **WHEN** a binary invokes the runtime with a handler and a listen address
- **THEN** the server binds the address and serves HTTP requests

#### Scenario: Graceful shutdown on signal

- **WHEN** the process receives SIGINT or SIGTERM
- **THEN** the server stops accepting new connections, waits for in-flight requests up to the shutdown timeout, and exits without a non-zero error for a clean shutdown

### Requirement: Liveness health endpoint

The service runtime SHALL expose a liveness health endpoint that reports whether the process is up. The endpoint SHALL respond successfully without depending on any downstream store or domain logic.

#### Scenario: Health check succeeds

- **WHEN** a client sends `GET /healthz`
- **THEN** the runtime responds with HTTP 200 and a body indicating the service is alive

### Requirement: Baseline middleware

The service runtime SHALL wrap handlers with baseline middleware applied to every request: assignment of a request ID, panic recovery, and request logging.

#### Scenario: Request ID assigned

- **WHEN** any request is handled
- **THEN** a request ID is generated (or taken from an inbound correlation header) and included on the response and in log entries for that request

#### Scenario: Panic recovery

- **WHEN** a handler panics while processing a request
- **THEN** the runtime recovers, logs the failure with the request ID, and returns a 500 response using the consistent error-response shape instead of dropping the connection

#### Scenario: Request logged

- **WHEN** a request completes
- **THEN** the runtime emits a structured log entry including method, path, status code, duration, and request ID

### Requirement: Consistent error-response shape

The service runtime SHALL provide a single helper for writing error responses so that all error responses across services share the same JSON shape and set the appropriate status code.

#### Scenario: Error written in standard shape

- **WHEN** a handler or middleware writes an error response
- **THEN** the body is JSON containing at least a machine-readable code and a human-readable message, and the HTTP status code matches the error

### Requirement: Environment-based configuration

The service runtime SHALL read configuration from the environment with sane defaults that can be overridden per environment. Configuration SHALL include at least the listen address, log level, and shutdown timeout.

#### Scenario: Defaults applied when unset

- **WHEN** no configuration environment variables are set
- **THEN** the runtime starts using documented default values

#### Scenario: Environment overrides defaults

- **WHEN** a configuration environment variable is set
- **THEN** the runtime uses the provided value instead of the default

### Requirement: Structured logging

The service runtime SHALL provide structured logging (via `log/slog`) at a configurable level, used by middleware and available to handlers.

#### Scenario: Log level honored

- **WHEN** the configured log level is set to a given threshold
- **THEN** log entries below that threshold are suppressed and entries at or above it are emitted
