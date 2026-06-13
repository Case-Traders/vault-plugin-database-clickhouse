# Installation

## Build the plugin binary

```bash
make build-linux-amd64   # bin/linux-amd64/clickhouse-database-plugin
make build-linux-arm64
make build-darwin-arm64
```

## Register with Vault

1. Copy the binary into Vault's plugin directory (often `/usr/local/libexec/vault/`).
2. Compute SHA256:

```bash
sha256sum clickhouse-database-plugin
```

3. Register:

```bash
vault login
vault plugin register -sha256=<sha256> database clickhouse-database-plugin
```

4. Configure a database secrets mount and connection as in [Configuration](configuration.md).

## Kubernetes (Helm values sketch)

Set `plugin_directory` in Vault server config:

```yaml
server:
  standalone:
    config: |
      plugin_directory = "/usr/local/libexec/vault/"
```

Download the release binary in an init container and mount it into that path. Example using the GitHub release API:

```yaml
extraInitContainers:
  - name: clickhouse-plugin
    image: alpine
    command: [sh, -c]
    args:
      - |
        set -e
        PLATFORM=linux-amd64
        VERSION=$(wget -qO- "https://api.github.com/repos/Case-Traders/vault-plugin-database-clickhouse/releases/latest" | grep '"tag_name"' | cut -d '"' -f 4)
        wget "https://github.com/Case-Traders/vault-plugin-database-clickhouse/releases/download/${VERSION}/clickhouse-database-plugin-${PLATFORM}.zip" -O plugin.zip
        unzip plugin.zip
        mv "${PLATFORM}/clickhouse-database-plugin" /plugins/clickhouse-database-plugin
        chmod +x /plugins/clickhouse-database-plugin
    volumeMounts:
      - name: plugins
        mountPath: /plugins
```

Mount `/plugins` at `/usr/local/libexec/vault` on the Vault container.
