# vault-plugin-database-clickhouse

External Vault database plugin for ClickHouse. Role statements in Vault drive user DDL; use `{{cluster}}` in SQL for `ON CLUSTER` on clusters.

![Architecture diagram](diagram.svg)

Full docs: [case-traders.github.io/vault-plugin-database-clickhouse](https://case-traders.github.io/vault-plugin-database-clickhouse/)

## Install

1. Binary from [GitHub Releases](https://github.com/Case-Traders/vault-plugin-database-clickhouse/releases) or `make build-linux-amd64`
2. `vault plugin register -sha256=… database clickhouse-database-plugin`
3. Database secrets mount + connection config

→ [Installation guide](docs/guides/installation.md) (K8s init container, paths)  
→ [Configuration guide](docs/guides/configuration.md) (keys, placeholders, Terraform)

## Development

```bash
devenv shell
ch-build
ch-ci            # Coq proofs + Rapid property tests
ch-integration   # + ClickHouse testcontainers (Docker; Podman: devenv up first)
```

→ [Development](docs/guides/development.md) · [Correctness](docs/guides/correctness.md)

## Acknowledgments

Thanks to [everythings-gonna-be-alright/vault-plugin-database-clickhouse](https://github.com/everythings-gonna-be-alright/vault-plugin-database-clickhouse) for the original plugin.

This is a Case-Traders rewrite, not a continuation of upstream releases.
