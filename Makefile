.PHONY: test test-race test-integration lint coverage check

test:
	go test ./...

test-race:
	go test -race ./...

test-integration:
	go test -tags=integration ./...

lint:
	go vet ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

check: lint test-race
