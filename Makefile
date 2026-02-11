BIN := cc-flavors

.PHONY: all
all: fmt lint test build

.PHONY: build
build:
	goreleaser build --single-target --snapshot --output $(BIN)

.PHONY: test
test:
	go test ./...

.PHONY: fmt
fmt:
	gofmt -w .

.PHONY: lint
lint:
	go vet ./...

.PHONY: clean
clean:
	rm -rf $(BIN) dist

.PHONY: release
release:
	goreleaser release --clean
