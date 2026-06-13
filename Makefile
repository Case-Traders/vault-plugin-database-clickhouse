PERENNIAL_PIN ?= c15c19774d4394959ae1e9ee85e5852df00046e7
PERENNIAL_WORK ?= $(CURDIR)/.cache/perennial
PERENNIAL_ROOT ?= $(PERENNIAL_WORK)
OPAMROOT ?= $(CURDIR)/.opam
OPAM_SWITCH ?= perennial-proof

GO ?= go
GOTOOLCHAIN ?= go1.25.11
export GOTOOLCHAIN

GO_TEST_FLAGS ?=
INTEGRATION_FLAGS ?= -tags=integration
GOMARKDOC ?= gomarkdoc

GOOSE_PKGS := ./internal/stmt ./internal/cluster/choose ./internal/txexec ./internal/vars ./internal/stmts ./internal/validate ./internal/deletepath
GOOSE_CODE := proof/goose/code
GOOSE_GEN := proof/goose/generatedproof
GOBIN ?= $(CURDIR)/.go/bin
GOPATH ?= $(CURDIR)/.go
GOMODCACHE ?= $(CURDIR)/.go/pkg/mod
export PATH := $(GOBIN):$(PATH)
export GOPATH GOMODCACHE OPAMROOT

.PHONY: default build build-linux-amd64 build-linux-arm64 build-darwin-arm64 \
	goose-tools goose proof-setup proof proof-local test test-integration \
	ci ci-integration docs-api docs docs-serve

default: build-linux-arm64

build-linux-amd64:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -o bin/linux-amd64/clickhouse-database-plugin ./clickhouse-database-plugin

build-linux-arm64:
	@CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GO) build -o bin/linux-arm64/clickhouse-database-plugin ./clickhouse-database-plugin

build-darwin-arm64:
	@CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GO) build -o bin/darwin-arm64/clickhouse-database-plugin ./clickhouse-database-plugin

goose-tools:
	@mkdir -p $(GOBIN) $(GOMODCACHE)
	@GOTOOLCHAIN=go1.26.0 GOBIN=$(GOBIN) GOMODCACHE=$(GOMODCACHE) \
		$(GO) install \
		github.com/mit-pdos/perennial/goose/cmd/goose@$(PERENNIAL_PIN) \
		github.com/mit-pdos/perennial/goose/cmd/proofgen@$(PERENNIAL_PIN) \
		github.com/mit-pdos/perennial/goose/cmd/proof-setup@$(PERENNIAL_PIN)

goose: goose-tools
	@rm -rf $(GOOSE_CODE)/vault_plugin_database_clickhouse $(GOOSE_GEN)/vault_plugin_database_clickhouse
	@for pkg in $(GOOSE_PKGS); do \
		goose -dir $(CURDIR) -out $(GOOSE_CODE) $$pkg; \
		proofgen -dir $(CURDIR) -out $(GOOSE_GEN) -configdir $(GOOSE_CODE) $$pkg; \
	done

proof/_CoqProject: proof/_CoqProject.in
	@sed \
		-e 's|@PERENNIAL_ROOT@|$(PERENNIAL_ROOT)|g' \
		-e 's|@PLUGIN_ROOT@|$(CURDIR)|g' \
		proof/_CoqProject.in > proof/_CoqProject

proof-setup:
	@if command -v ch-proof-setup >/dev/null 2>&1; then ch-proof-setup; else devenv shell -- ch-proof-setup; fi

proof-local: proof-setup proof/_CoqProject
	@eval "$$(OPAMROOT=$(OPAMROOT) opam env --switch=$(OPAM_SWITCH))" && \
	cd proof && \
	rocq makefile -f _CoqProject -o Makefile.coq && \
	$(MAKE) -f Makefile.coq all

proof: goose proof-local

test:
	$(GO) test $(GO_TEST_FLAGS) ./...

test-integration:
	$(GO) test $(INTEGRATION_FLAGS) $(GO_TEST_FLAGS) ./...

ci: goose proof

ci-integration: ci test-integration

docs-api:
	@rm -rf docs/api && mkdir -p docs/api
	$(GOMARKDOC) -o 'docs/api/{{if eq .Dir "."}}clickhouse{{else}}{{.Dir}}{{end}}.md' \
		. ./clickhouse-database-plugin ./internal/cluster ./internal/cluster/choose ./internal/stmt \
		./internal/txexec ./internal/user ./internal/vars ./internal/stmts \
		./internal/validate ./internal/deletepath ./testutil

docs: docs-api
	@cp diagram.svg docs/diagram.svg
	mkdocs build

docs-serve: docs-api
	@cp diagram.svg docs/diagram.svg
	mkdocs serve
