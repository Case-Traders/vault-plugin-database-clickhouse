package cluster

import "errors"

var (
	// ErrEmptyCluster is returned when system.clusters has no rows and cluster is unset.
	ErrEmptyCluster = errors.New("no ClickHouse cluster discovered")
	// ErrAmbiguousCluster is returned when multiple clusters exist and cluster is unset.
	ErrAmbiguousCluster = errors.New("multiple ClickHouse clusters discovered; set cluster in Vault database config")
)
