# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- API endpoint infrastructure: `redirectd` and `shortend` service binaries, each exposing a placeholder route (`GET /{code}`, `POST /api/shorten`) over a shared HTTP runtime (`internal/platform`) providing graceful lifecycle, request-ID/panic-recovery/request-logging middleware, a liveness health endpoint (`GET /healthz`), a consistent JSON error-response shape, structured logging, and environment-based configuration.
- Makefile targets (`build`, `run-redirectd`, `run-shortend`, `test`, `lint`) mirroring CI checks.
