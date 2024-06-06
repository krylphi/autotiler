GO ?= GOPRIVATE=github.com/Rent-Set go
GOLINT ?= golint
# with caching
#GOLANGCILINT ?= docker run --rm -v $(shell pwd):/app -v ~/.cache/golangci-lint/v1.57.2:/root/.cache -w /app golangci/golangci-lint:v1.57.2 golangci-lint
# no caching
GOLANGCILINT ?= docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v1.57.2 golangci-lint
GOSEC ?= gosec
GOPATH ?= $(shell go env GOPATH)
VERSION ?=$(shell git describe --tags --always)
PACKAGES = $(shell go list -f {{.Dir}} ./... | grep -v /vendor/ | grep -v /proto )
DATE = $(shell date -R)

BIN_NAME=autotiler

.PHONY: build
build:
	go build -o ./out/$(BIN_NAME) .

.PHONY: tidy
tidy:
	$(GO) mod tidy

.PHONY: mod-download
mod-download:
	$(GO) mod download

.PHONY: lint
lint: ## Lint the codebase.
	$(GO) vet ${PACKAGES}
	$(GOLANGCILINT) run -v

.PHONY: gosec
gosec:
	gosec -fmt=json -out=gosec.out.local.json ./...

.PHONY: src-fmt
src-fmt:
	gofmt -s -w ${PACKAGES}
	gci write --skip-generated --skip-vendor --section Standard --section Default --section "Prefix(github.com/krylphi)" --section "Prefix(github.com/krylphi/$(BIN_NAME))" ${PACKAGES}

.PHONY: unpack-example
unpack-example:
	go run . ./examples/2x3_packed.png output.local.png\

