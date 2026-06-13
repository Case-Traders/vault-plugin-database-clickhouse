// Package clickhouse implements a HashiCorp Vault database secrets plugin for ClickHouse.
//
// The plugin creates and rotates database users from Vault role statements.
// Cluster names for ON CLUSTER DDL come from Vault config or from system.clusters.
package clickhouse
