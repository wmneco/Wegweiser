# Contributing to Wegweiser

Thank you for your interest in contributing!

## Getting Started

1. Fork the repository and clone it locally.
2. Create a feature branch: `git checkout -b feat/my-feature`
3. Make your changes and add tests where applicable.
4. Run the tests: `go test ./...`
5. Run static analysis: `go vet ./...`
6. Commit your changes following [Conventional Commits](https://www.conventionalcommits.org/).
7. Open a Pull Request against `main`.

## Code Style

- Follow standard Go conventions (`gofmt`, `goimports`).
- Keep functions small and focused.
- Public APIs must be documented with Go doc comments.

## Reporting Issues

Please use the GitHub issue tracker. Include:
- A clear description of the problem.
- Steps to reproduce.
- Expected vs. actual behavior.
- Go version and OS.

## Pull Requests

- One concern per PR.
- Update `CHANGELOG.md` under `## [Unreleased]`.
- Ensure all CI checks pass before requesting review.
