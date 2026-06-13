package cluster

import (
	"context"
	"database/sql"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
	P "github.com/IBM/fp-go/v2/predicate"

	"vault-plugin-database-clickhouse/internal/cluster/choose"
)

// Resolve returns the cluster name using configured override or discovery.
func Resolve(ctx context.Context, db *sql.DB, configured string) (string, error) {
	return E.Uneitherize0(func() E.Either[error, string] {
		return F.Pipe2(
			configured,
			O.FromPredicate(P.IsNonZero[string]()),
			O.Fold(
				func() E.Either[error, string] { return resolveDiscovered(ctx, db) },
				func(cfg string) E.Either[error, string] {
					return F.Pipe1(cfg, E.Eitherize1(func(c string) (string, error) {
						return choose.ChooseCluster(c, nil)
					}))
				},
			),
		)
	})()
}

func resolveDiscovered(ctx context.Context, db *sql.DB) E.Either[error, string] {
	return F.Pipe1(
		E.Eitherize2(Discover)(ctx, db),
		E.Chain(pickDiscoveredCluster),
	)
}

func pickDiscoveredCluster(discovered []string) E.Either[error, string] {
	return F.Pipe1(discovered, E.Eitherize1(func(ds []string) (string, error) {
		return choose.ChooseCluster("", ds)
	}))
}
