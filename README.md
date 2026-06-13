# vault-plugin-database-clickhouse

Vault database plugin for ClickHouse. Creates users from Vault role statements.

Supports single-node and clustered ClickHouse (`ON CLUSTER '{{cluster}}'` in your statements).

![Diagram](diagram.svg)

Documentation: https://case-traders.github.io/vault-plugin-database-clickhouse/

## Example

```hcl
resource "vault_database_secret_backend_role" "clickhouse_analytics" {
  backend = vault_mount.clickhouse.path
  name    = "clickhouse_analytics"
  db_name = "my_clickhouse"
  creation_statements = [
    "CREATE USER '{{name}}' IDENTIFIED WITH sha256_password BY '{{password}}' ON CLUSTER '{{cluster}}';",
    "GRANT ON CLUSTER '{{cluster}}' analytics TO '{{name}}';",
  ]
  default_ttl = 2593000
  max_ttl     = 2593000
}
```

More keys and placeholders: [Configuration guide](docs/guides/configuration.md).

## Install

Download a binary from [GitHub Releases](https://github.com/Case-Traders/vault-plugin-database-clickhouse/releases) or run `make build-linux-amd64`.

Copy it into Vault's plugin directory, compute SHA256, register:

```bash
sha256sum clickhouse-database-plugin
vault login
vault plugin register -sha256=<sha256> database clickhouse-database-plugin
```

Then configure a database secrets mount and connection. See [Installation guide](docs/guides/installation.md) for paths and Kubernetes init-container setup.

## Build and test

```bash
devenv shell   # optional
make build-linux-amd64
make ci
make ci-integration   # Docker or Podman; see development guide
```

[Development guide](docs/guides/development.md). [Correctness guide](docs/guides/correctness.md) (Goose, Perennial, integration tests).

## Acknowledgments

Thanks to [everythings-gonna-be-alright/vault-plugin-database-clickhouse](https://github.com/everythings-gonna-be-alright/vault-plugin-database-clickhouse) for the original plugin. Their commits are still in this repo's history.
