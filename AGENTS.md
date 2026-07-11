# Agent Instructions for Wegweiser

## Project Overview

Wegweiser is a Go project located at `github.com/wmneco/wegweiser`.

## Repository Layout

```
docs/           # Documentation (Diátaxis structure)
  explanation/  # Background and conceptual explanations
  how-to/       # Goal-oriented guides
  reference/    # API / technical reference
  tutorials/    # Step-by-step learning material
  assets/       # Images and other static assets
```

## Conventions

- Language: Go — follow idiomatic Go style (`gofmt`, `goimports`).
- Commit style: [Conventional Commits](https://www.conventionalcommits.org/) (`feat:`, `fix:`, `docs:`, `chore:`, etc.).
- All user-facing changes must be recorded in `CHANGELOG.md` under `## [Unreleased]`.
- Documentation follows the [Diátaxis](https://diataxis.fr/) framework.

## Do Not

- Commit directly to `main`.
- Skip updating `CHANGELOG.md` for user-facing changes.
- Add dependencies without discussing the trade-offs first.
