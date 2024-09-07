.PHONY: ci golang-imports golang-imports-install install
.PHONY: lint-install lint mock-clean mock-install mock-gen
.PHONY: test-html-output test

LINT_VERSION=v1.60.3
MOQ_VERSION=v0.3.4
SHELL := /bin/bash

ci: golang-imports lint mock-gen test

golang-imports:
	goimports -w .

golang-imports-install:
	go install golang.org/x/tools/cmd/goimports@latest

install:
	go install

lint: mock-clean golang-imports
	golangci-lint run --config golangci.yaml --timeout 10m

lint-install:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin $(LINT_VERSION)

mock-clean:
	rm -rf mocks
	mkdir -p mocks

mock-gen: mock-clean
	rm -rf mocks
	mkdir -p mocks
	go generate ./internal/...

mock-install:
	go install github.com/matryer/moq@$(MOQ_VERSION)

test:
	go test -cover ./...

test-html-output:
	go test -coverprofile=c.out ./... && go tool cover -html=c.out && rm -f c.out
