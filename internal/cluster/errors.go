package cluster

import "vault-plugin-database-clickhouse/internal/cluster/choose"

var (
	// ErrEmptyCluster is returned when system.clusters has no rows and cluster is unset.
	ErrEmptyCluster = choose.ErrEmptyCluster
	// ErrAmbiguousCluster is returned when multiple clusters exist and cluster is unset.
	ErrAmbiguousCluster = choose.ErrAmbiguousCluster
)
