GO ?= go
GOTOOLCHAIN ?= go1.25.11
export GOTOOLCHAIN

GO_TEST_FLAGS ?=
INTEGRATION_FLAGS ?= -tags=integration
GOMARKDOC ?= gomarkdoc

.PHONY: default build build-linux-amd64 build-linux-arm64 build-darwin-arm64 \
	proof test test-integration ci ci-integration docs-api docs docs-serve

default: build-linux-arm64

build-linux-amd64:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -o bin/linux-amd64/clickhouse-database-plugin ./clickhouse-database-plugin

build-linux-arm64:
	@CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GO) build -o bin/linux-arm64/clickhouse-database-plugin ./clickhouse-database-plugin

build-darwin-arm64:
	@CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GO) build -o bin/darwin-arm64/clickhouse-database-plugin ./clickhouse-database-plugin

proof:
	@cd proof && coq_makefile -f _CoqProject -o Makefile.coq && $(MAKE) -f Makefile.coq clean && $(MAKE) -f Makefile.coq

test:
	$(GO) test $(GO_TEST_FLAGS) ./...

test-integration:
	$(GO) test $(INTEGRATION_FLAGS) $(GO_TEST_FLAGS) ./...

ci: proof test

ci-integration: ci test-integration

docs-api:
	@rm -rf docs/api && mkdir -p docs/api
	$(GOMARKDOC) -o 'docs/api/{{if eq .Dir "."}}clickhouse{{else}}{{.Dir}}{{end}}.md' \
		. ./clickhouse-database-plugin ./internal/cluster ./internal/stmt \
		./internal/txexec ./internal/user ./internal/vars ./testutil

docs: docs-api
	@cp diagram.svg docs/diagram.svg
	mkdocs build

docs-serve: docs-api
	@cp diagram.svg docs/diagram.svg
	mkdocs serve
