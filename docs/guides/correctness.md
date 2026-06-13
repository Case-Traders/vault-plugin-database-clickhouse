# Verification

How this plugin is checked before release.

## Layers

| Layer | Tool | Command | What it checks |
| ----- | ---- | ------- | -------------- |
| Formal | Coq in `proof/` | `make proof` | Theorems on statement, cluster, execution, and template specs |
| Properties | Rapid | `go test ./...` | Go code matches specs and golden JSON vectors |
| Integration | testcontainers | `go test -tags=integration ./...` | Real ClickHouse DDL and Vault dbplugin RPC |

Pure packages in `internal/stmt`, `internal/cluster`, `internal/txexec`, and `internal/vars` use Rapid in `*_rapid_test.go`. There are no parallel table-driven unit files for the same logic.

Coq modules map to Go as follows:

| Coq | Go |
| --- | --- |
| `proof/Stmt.v` | `internal/stmt.NormalizeCommands` |
| `proof/Cluster.v` | `internal/cluster.ChooseCluster` |
| `proof/FirstError.v` | `internal/txexec.FirstError` |
| `proof/Vars.v` | `internal/vars.HasRequiredKeys` |

## Coq theorems

### Statement normalization (`proof/Stmt.v`)

| Theorem | Statement |
| ------- | --------- |
| `normalize_idempotent` | `normalize (normalize cs) = normalize cs` |

Golden vectors live in `proof/testvectors/stmt_normalize.json` and load in Rapid tests.

### Cluster selection (`proof/Cluster.v`)

| Theorem | Statement |
| ------- | --------- |
| `choose_configured` | Configured name is always chosen |
| `choose_single_discovery` | Exactly one discovered cluster is used |
| `choose_empty_error` | No clusters returns error |

### Ordered execution (`proof/FirstError.v`)

| Lemma | Statement |
| ----- | --------- |
| `first_error_stops` | Prefix of length *i* exists before failure index *i* |

### Template placeholders (`proof/Vars.v`)

| Lemma | Statement |
| ----- | --------- |
| `new_user_keys_complete` | Default new-user placeholders cover required keys |

## Rapid property tests

| File | Property |
| ---- | -------- |
| `internal/stmt/normalize_rapid_test.go` | Idempotence and Coq golden JSON |
| `internal/cluster/choose_rapid_test.go` | Configured, empty, single, ambiguous clusters |
| `internal/txexec/exec_rapid_test.go` | Fail-fast and error preservation |
| `internal/vars/builder_rapid_test.go` | Required keys per operation |

## Integration tests

| Test | Files |
| ---- | ----- |
| Initialize | `clickhouse_integration_test.go` |
| NewUser, UpdateUser, DeleteUser | `clickhouse_integration_test.go`, `clickhouse-database-plugin/plugin_integration_test.go` |
| UpdateUser missing user | `clickhouse_integration_test.go` |
| User exists | `internal/user/exists_integration_test.go` |
| Cluster config override | `clickhouse_integration_test.go` |

## CI commands

```bash
make ci              # proof + Rapid
make ci-integration  # proof + Rapid + testcontainers
```
