.PHONY: build test

OS   := $(shell uname -s | tr A-Z a-z)
ARCH := $(shell uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/')
OUT  := plugin/bin/coderay-skeleton-$(OS)-$(ARCH)

# Build pre-built binary for the current platform and place it in bin/
build:
	CGO_ENABLED=1 go build -o $(OUT) ./cmd/coderay-skeleton
	@echo "built $(OUT)"

test:
	CGO_ENABLED=1 go test ./...
