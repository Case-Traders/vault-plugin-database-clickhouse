// Package cluster resolves ClickHouse cluster names for ON CLUSTER DDL.
//
// When Vault config sets cluster, that name is used. Otherwise names are read
// from system.clusters and ChooseCluster applies single-or-error rules.
package cluster
