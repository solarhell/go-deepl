test:
	go test -short -race ./...

e2e-test:
	./scripts/e2e-test $(authKey)

lint:
	go mod tidy -v
	go vet ./...
	golangci-lint run

.PHONY: test e2e-test lint
