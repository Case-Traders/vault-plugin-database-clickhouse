package cluster

import (
	"context"
	"database/sql"
	"fmt"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
)

const discoverQuery = `SELECT DISTINCT cluster FROM system.clusters ORDER BY cluster`

var mapQueryError = E.MapLeft[*sql.Rows, error, error](wrapQueryError)
var mapScanError = E.MapLeft[string, error, error](wrapScanError)
var mapIterateError = E.MapLeft[[]string, error, error](wrapIterateError)

// Discover returns distinct cluster names from system.clusters.
func Discover(ctx context.Context, db *sql.DB) ([]string, error) {
	return E.Uneitherize0(func() E.Either[error, []string] {
		return withQueryRows(ctx, db)(collectClusterNames)
	})()
}

var withQueryRows = func(ctx context.Context, db *sql.DB) func(func(*sql.Rows) E.Either[error, []string]) E.Either[error, []string] {
	return E.WithResource[[]string, error, *sql.Rows, any](
		func() E.Either[error, *sql.Rows] {
			return mapQueryError(E.TryCatchError(db.QueryContext(ctx, discoverQuery)))
		},
		releaseRows,
	)
}

func releaseRows(rows *sql.Rows) E.Either[error, any] {
	_ = rows.Close()
	return E.Of[error, any](nil)
}

func collectClusterNames(rows *sql.Rows) E.Either[error, []string] {
	return foldClusterRows(rows, nil)
}

func foldClusterRows(rows *sql.Rows, acc []string) E.Either[error, []string] {
	return F.Pipe2(
		rows.Next(),
		O.FromPredicate(func(b bool) bool { return b }),
		O.Fold(
			func() E.Either[error, []string] { return finishClusterRows(rows, acc) },
			func(bool) E.Either[error, []string] {
				return E.Chain(func(name string) E.Either[error, []string] {
					return foldClusterRows(rows, append(acc, name))
				})(scanClusterRow(rows))
			},
		),
	)
}

func scanClusterRow(rows *sql.Rows) E.Either[error, string] {
	return mapScanError(E.Eitherize1(func(r *sql.Rows) (string, error) {
		var name string
		return name, r.Scan(&name)
	})(rows))
}

func finishClusterRows(rows *sql.Rows, acc []string) E.Either[error, []string] {
	return mapIterateError(E.TryCatchError(acc, rows.Err()))
}

func wrapQueryError(err error) error {
	return fmt.Errorf("query system.clusters: %w", err)
}

func wrapScanError(err error) error {
	return fmt.Errorf("scan cluster: %w", err)
}

func wrapIterateError(err error) error {
	return fmt.Errorf("iterate clusters: %w", err)
}
