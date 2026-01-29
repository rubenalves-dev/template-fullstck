.PHONY: test test-coverage lint build

test:
	go test ./... -v

test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

lint:
	golangci-lint run

build:
	go build -o bin/api cmd/api/main.go

watch:
	@if command -v air > /dev/null; then \
	    air; \
	else \
	    echo "Air is not installed. Please install it with: go install github.com/air-verse/air@latest"; \
	fi
