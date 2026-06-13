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

GOOSE_PKGS := ./internal/stmt ./internal/cluster/choose ./internal/firsterror ./internal/vars ./internal/validate ./internal/deletepath
GOOSE_CODE := proof/goose/code
GOOSE_GEN := proof/goose/generatedproof
GOOSE_STUBS := $(CURDIR)/proof/goose/proof-stubs
# Minimal Perennial .vo files required by plugin generatedproof modules.
PERENNIAL_PROOF_VOS := \
	new/proof/proof_prelude.vo \
	new/proof/strings.vo \
	new/proof/slices_proof/slices_init.vo \
	new/proof/errors.vo \
	new/proof/io.vo \
	new/proof/fmt.vo \
	new/proof/sort_proof/sort_init.vo \
	new/golang/theory.vo \
	new/code/slices.vo \
	new/code/strings.vo \
	new/code/sort.vo \
	new/code/errors.vo \
	new/code/fmt.vo \
	new/generatedproof/slices.vo \
	new/generatedproof/strings.vo \
	new/generatedproof/errors.vo \
	new/generatedproof/fmt.vo \
	new/generatedproof/sort.vo
GOBIN ?= $(CURDIR)/.go/bin
GOPATH ?= $(CURDIR)/.go
GOMODCACHE ?= $(CURDIR)/.go/pkg/mod
export PATH := $(GOBIN):$(PATH)
export GOPATH GOMODCACHE OPAMROOT

.PHONY: default build build-linux-amd64 build-linux-arm64 build-darwin-arm64 \
	goose-tools goose proof-deps proof-stubs proof-setup proof proof-local test test-integration \
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

proof-deps:
	@if test -f "$(PERENNIAL_ROOT)/new/proof/proof_prelude.vo"; then \
		echo "proof-deps: Perennial ready ($(PERENNIAL_ROOT))"; \
	elif command -v ch-proof-setup >/dev/null 2>&1; then \
		ch-proof-setup; \
	elif command -v devenv >/dev/null 2>&1; then \
		devenv shell -- ch-proof-setup; \
	else \
		echo "proof-deps: missing $(PERENNIAL_ROOT)/new/proof/proof_prelude.vo"; \
		echo "proof-deps: install Rocq/opam deps and build proof_prelude.vo, or run ch-proof-setup in devenv shell"; \
		exit 1; \
	fi
	@$(MAKE) -C $(PERENNIAL_ROOT) -j$$(nproc 2>/dev/null || echo 2) $(PERENNIAL_PROOF_VOS)

proof-stubs: goose proof-deps
	@for pkg in $(GOOSE_PKGS); do \
		proof-setup -dir $(CURDIR) -out $(GOOSE_STUBS) $$pkg; \
	done
	@for f in $(GOOSE_STUBS)/vault_plugin_database_clickhouse/internal/*.v \
		$(GOOSE_STUBS)/vault_plugin_database_clickhouse/internal/cluster/*.v; do \
		sed -i \
			-e '/Require Import maps\./d' \
			-e 's/ `{!globalsGS Σ} {go_ctx: GoContext}//g' \
			-e 's/: IsPkgInit \([a-z_]*\) :=/: IsPkgInit (iProp Σ) \1 :=/g' \
			-e 's/: GetIsPkgInitWf \([a-z_]*\) :=/: GetIsPkgInitWf (iProp Σ) \1 :=/g' \
			"$$f"; \
		if ! grep -q 'Require Export proof_prelude' "$$f"; then \
			sed -i '1i From New.proof Require Export proof_prelude.' "$$f"; \
		fi; \
	done

# Deprecated alias.
proof-setup: proof-deps

proof-local: proof-deps proof/_CoqProject
	@if [ -d "$(OPAMROOT)/$(OPAM_SWITCH)" ]; then \
		eval "$$(OPAMROOT=$(OPAMROOT) opam env --switch=$(OPAM_SWITCH))"; \
	else \
		eval "$$(opam env)"; \
	fi && \
	cd proof && \
	rocq makefile -f _CoqProject -o Makefile.coq && \
	$(MAKE) -f Makefile.coq clean && \
	$(MAKE) -f Makefile.coq all && \
	for f in deletepath_forall stmt_forall choose_forall firsterror_forall vars_forall validate_forall; do \
		test -f clickhouse/$$f.vo || (echo "proof-local: missing clickhouse/$$f.vo" && exit 1); \
	done && \
	if grep -rE 'Admitted|Abort\.' clickhouse/*.v; then \
		echo "proof-local: found Admitted or Abort in proof/clickhouse/*.v" && exit 1; \
	fi

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
		./internal/firsterror ./internal/txexec ./internal/user ./internal/vars \
		./internal/validate ./internal/deletepath ./testutil

docs: docs-api
	@cp diagram.svg docs/diagram.svg
	mkdocs build

docs-serve: docs-api
	@cp diagram.svg docs/diagram.svg
	mkdocs serve
