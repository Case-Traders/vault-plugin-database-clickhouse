# Development

Local setup for building, verifying, and documenting the plugin.

## devenv

```bash
devenv shell
```

Provides Go 1.25, Goose tools, **opam from nixpkgs** (Rocq/OCaml via project opam switch), and Docker or Podman for testcontainers.

| Command | Purpose |
| ------- | ------- |
| `make build-linux-amd64` | Build plugin binary |
| `make goose` | Translate provable Go packages to Coq |
| `ch-proof-setup` | One-time opam switch + Perennial + Rocq (in shell) |
| `make proof-setup` | Same as `ch-proof-setup` (works outside shell too) |
| `ch-proof` or `make proof` | Compile Perennial Iris proofs |
| `make ci` | `goose` + `proof` |
| `make test-integration` | testcontainers integration tests |
| `make ci-integration` | `ci` + integration tests |
| `make docs` | API + MkDocs site |

Shell scripts: `ch-build`, `ch-proof-setup`, `ch-proof`, `ch-ci`, …

## Environment

GMP and libc headers come from nix (`env.PKG_CONFIG_PATH`, `env.CPATH` in devenv). Run `ch-proof-setup` once inside `devenv shell`.

**First run is slow:** after Rocq installs, opam compiles **rocq-iris** (Iris) from pinned git sources. That step often takes **5-10 minutes** with minimal terminal output while `make`/`rocq` runs in the background. High CPU usage means it is working, not hung.

## Containers

Integration tests need Docker or Podman. devenv detects `/var/run/docker.sock`; otherwise it starts a Podman API socket (`devenv up`).

## Docs

```bash
make docs-serve
```

See [Correctness guide](correctness.md) for Goose + Perennial verification.
