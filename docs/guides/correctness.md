# Verification

How this plugin is checked before release.

## Layers

| Layer | Tool | Command | What it checks |
| ----- | ---- | ------- | -------------- |
| Goose + Perennial | [Goose](https://github.com/goose-lang/goose) + [Perennial](https://github.com/mit-pdos/perennial) | `make proof` | Go source is translated to Coq; Iris WP proofs on `test*()` semantics functions |
| Integration | testcontainers | `go test -tags=integration ./...` | Real ClickHouse DDL and Vault dbplugin RPC |

Pure packages under `internal/` use imperative Go (Goose-compatible). `test*() bool` functions in `semantics.go` files are translated by Goose and proved in Perennial under [`proof/clickhouse/`](../../proof/clickhouse/).

## Goose pipeline

```bash
make goose       # goose + proofgen → proof/goose/{code,generatedproof}/
make proof       # compile Perennial proofs (Rocq 9.x via devenv nix input)
make ci          # goose + proof
make ci-integration  # ci + integration tests
```

Generated Coq is committed as gold in `proof/goose/`. Perennial is pinned by `PERENNIAL_PIN` and provided via `PERENNIAL_ROOT` (devenv nix input or a local checkout).

## Proof layout

| Path | Role |
| ---- | ---- |
| `proof/goose/code/` | Goose translation (`New.code.*`) |
| `proof/goose/generatedproof/` | Proofgen wrappers (`New.generatedproof.*`) |
| `proof/clickhouse/init.v` | Shared `test_fun_ok` and tactics |
| `proof/clickhouse/*_forall.v` | lemmas |

## Semantics proofs

| Go package | Lemma file | Property |
| ---------- | ---------- | -------- |
| `internal/stmt` | `stmt_forall.v` | Idempotence|
| `internal/cluster/choose` | `choose_forall.v` | Configured, empty, single, ambiguous clusters |
| `internal/txexec` | `txexec_forall.v` | Fail-fast execution and error preservation |
| `internal/vars` | `vars_forall.v` | Template placeholders |
| `internal/stmts` | `stmts_forall.v` | Statement fallback |
| `internal/validate` | `validate_forall.v` | Request guards |
| `internal/deletepath` | `deletepath_forall.v` | DeleteUser routing |

## Integration tests

| Test | Files |
| ---- | ----- |
| Initialize | `clickhouse_integration_test.go` |
| NewUser, UpdateUser, DeleteUser | `clickhouse_integration_test.go`, `clickhouse-database-plugin/plugin_integration_test.go` |
| UpdateUser missing user | `clickhouse_integration_test.go` |
| User exists | `internal/user/exists_integration_test.go` |
| Cluster config override | `clickhouse_integration_test.go` |

## Toolchain setup

1. `devenv shell`
2. First time only: `ch-proof-setup` (or `make proof-setup`)
3. `make goose` after changing provable Go
4. `ch-proof`

`ch-proof-setup` installs Rocq via opam, Perennial dependencies, and builds `new/proof/proof_prelude.vo` in `.cache/perennial`. Re-run when `PERENNIAL_PIN` changes.

Optional fallback: set `PERENNIAL_ROOT` to your own Perennial checkout per [Perennial opam guide](https://github.com/mit-pdos/perennial/blob/master/docs/opam.md).

SQL I/O (`cluster.Discover`, `user.Exists`, Vault RPC) is not Goose-translated; integration tests cover it.
