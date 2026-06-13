# Configuration

Vault database secrets engine settings for this plugin.

## Keys

| Key | Required | Description |
| --- | -------- | ----------- |
| `connection_url` | yes | ClickHouse DSN (native or HTTP) |
| `username` / `password` | if URL is templated | Admin credentials used to run DDL |
| `cluster` | no | Cluster name for `ON CLUSTER` in default statements. When empty, the plugin reads `system.clusters`. Zero rows: error. One row: use that name. Multiple rows: error until you set `cluster` |
| `username_template` | no | Vault username template. Default generates `v-<display>-<role>-<random>-<time>` truncated to 63 characters |

## Statement placeholders

Vault replaces these in role statements:

| Placeholder | NewUser | Update password | Update expiration | DeleteUser |
| ----------- | ------- | --------------- | ----------------- | ---------- |
| `{{name}}` / `{{username}}` | yes | yes | yes | yes |
| `{{password}}` | yes | yes | no | no |
| `{{expiration}}` | yes | no | yes | no |
| `{{cluster}}` | yes | yes | yes | no |

Default password and expiration statements (when role statements are empty) use `ON CLUSTER '{{cluster}}'`.

## Terraform example

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

ClickHouse 24.8 integration tests use `IDENTIFIED WITH plaintext_password BY '{{password}}'` because the test image expects that form.
