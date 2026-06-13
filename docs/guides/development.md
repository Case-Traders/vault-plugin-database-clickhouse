# Development

## Environment

```bash
devenv shell
```

Provides Go 1.25, Coq 8.18, coq-lsp, and Docker or Podman for testcontainers.

## Make targets

| Target | Action |
| ------ | ------ |
| `make test` | Rapid property tests |
| `make proof` | Coq proofs in `proof/` |
| `make ci` | `proof` + `test` |
| `make test-integration` | ClickHouse testcontainers (needs container runtime) |
| `make ci-integration` | `ci` + integration tests |
| `make docs` | Build GitHub Pages site into `site/` |
| `make docs-serve` | Preview docs at http://127.0.0.1:8000 |
| `goreleaser release --snapshot --clean --skip=sign` | Local release dry run into `dist/` |

devenv scripts: `ch-test`, `ch-proof`, `ch-ci`, `ch-test-integration`, `ch-integration`, `ch-docs`, `ch-docs-serve`, `ch-release-snapshot`.

## Integration tests

testcontainers starts `clickhouse/clickhouse-server:24.8-alpine`. With Podman, run `devenv up` or enable `podman.socket` so `DOCKER_HOST` points at the user socket.

Plugin RPC tests live in `clickhouse-database-plugin/plugin_integration_test.go` and call `Run()` in a go-plugin subprocess.
