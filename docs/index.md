# vault-plugin-database-clickhouse

External Vault database plugin for ClickHouse. Role statements in Vault drive user DDL; use `{{cluster}}` in SQL for `ON CLUSTER` on clusters.

![Architecture](../diagram.svg)

## Guides

- [Configuration](guides/configuration.md): Vault database config keys and Terraform examples
- [Installation](guides/installation.md): plugin binary, registration, Kubernetes
- [Development](guides/development.md): devenv, make targets, tests
- [Correctness](guides/correctness.md): Coq, Rapid, integration tests

## API

Generated from Go doc comments. Run `make docs-api` locally, or see the **API** section in the site nav after `make docs`.
