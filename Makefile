.PHONY: build run-redirectd run-shortend test lint fmt vet

build:
	go build ./...

run-redirectd:
	go run ./cmd/redirectd

run-shortend:
	go run ./cmd/shortend

test:
	go test -race -coverprofile=coverage.out ./...

fmt:
	@unformatted=$$(go list -f '{{.Dir}}' ./... | xargs -r gofmt -l); \
	if [ -n "$$unformatted" ]; then \
		echo "The following files are not gofmt-formatted:"; \
		echo "$$unformatted"; \
		exit 1; \
	fi
	@unformatted=$$(go list -f '{{.Dir}}' ./... | xargs -r go tool goimports -l); \
	if [ -n "$$unformatted" ]; then \
		echo "The following files have unorganized imports:"; \
		echo "$$unformatted"; \
		exit 1; \
	fi

vet:
	go vet ./...

lint: fmt vet
