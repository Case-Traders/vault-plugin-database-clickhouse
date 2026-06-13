# Verification

The plugin is checked with Coq proofs on selected `internal/` packages and with Docker integration tests against ClickHouse.

## Verification layers

| Layer | Tool | Command | Covers |
| ----- | ---- | ------- | ------ |
| Formal | [Goose](https://github.com/goose-lang/goose) + [Perennial](https://github.com/mit-pdos/perennial) | `make proof` | `test*()` functions in `semantics.go`; Iris WP termination |
| Integration | testcontainers | `go test -tags=integration ./...` | SQL, Vault dbplugin RPC, plugin lifecycle |

## Proof layout

| Path | Contents |
| ---- | -------- |
| `proof/goose/code/` | Goose translation (`New.code.*`) |
| `proof/goose/generatedproof/` | proofgen wrappers (`New.generatedproof.*`) |
| `proof/clickhouse/init.v` | `test_fun_ok`, `clickhouse_semantics_auto` |
| `proof/clickhouse/*_forall.v` | One lemma per `test*()` function |

Six packages: `stmt`, `cluster/choose`, `firsterror`, `vars`, `validate`, `deletepath`. Gold under `proof/goose/` is committed. Perennial pin: `PERENNIAL_PIN`. Tree: `PERENNIAL_ROOT` (default `.cache/perennial`).

## Makefile targets

| Target | Action |
| ------ | ------ |
| `make goose` | Regenerate `proof/goose/{code,generatedproof}/` |
| `make proof-local` | Build Perennial deps, compile `proof/clickhouse/*.v`, require all six `*_forall.vo` |
| `make proof` | `goose`, then `proof-local` |
| `make ci` | `goose` + `proof` |
| `make ci-integration` | `ci` + `go test -tags=integration ./...` |

Cached `proof_prelude.vo` skips rebuilding Perennial only; plugin lemmas still compile each `proof-local` run.

## Semantics lemmas

Lemmas run the Goose translation with `clickhouse_semantics_auto`. No admitted WP layer for plugin code.

| Package | Lemma file | Proved in `semantics.go` | Not proved here |
| ------- | ---------- | ------------------------ | --------------- |
| `stmt` | `stmt_forall.v` | `StatementsOrDefault`, slice length | `NormalizeCommands`, `strings.Split` |
| `choose` | `choose_forall.v` | `configuredWins`, `discoveryRequired` | `ChooseCluster`, `strings.TrimSpace`, `sort` |
| `firsterror` | `firsterror_forall.v` | `nilStep` returns nil | `FirstError` with closures |
| `vars` | `vars_forall.v` | `OpNewUser`, `OpDeleteUser` values | `ForNewUser`, maps, `HasRequiredKeys` |
| `validate` | `validate_forall.v` | `hasUsername`, `needsCreationStatements` | `UpdateUser` with `fmt.Errorf` |
| `deletepath` | `deletepath_forall.v` | `UseCustomRevocation` on nil / one-element slice | — |

Full behaviour for cluster selection, validation errors, and template keys is in integration tests.

## Integration tests

| Area | Files |
| ---- | ----- |
| Plugin lifecycle, NewUser, UpdateUser, DeleteUser | `clickhouse_integration_test.go`, `clickhouse-database-plugin/plugin_integration_test.go` |
| User exists | `internal/user/exists_integration_test.go` |
| Not Goose-translated | `cluster.Discover`, `user.Exists`, Vault RPC |

## Setup

| Step | Command |
| ---- | ------- |
| Enter environment | `devenv shell` |
| First-time Rocq / Perennial | `ch-proof-setup` |
| After `semantics.go` or provable Go changes | `make goose` |
| Compile proofs | `make proof` or `ch-proof` |

Re-run `ch-proof-setup` when `PERENNIAL_PIN` changes. External Perennial: set `PERENNIAL_ROOT` ([opam guide](https://github.com/mit-pdos/perennial/blob/master/docs/opam.md)).
