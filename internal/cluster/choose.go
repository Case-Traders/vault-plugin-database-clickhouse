package cluster

import (
	"cmp"
	"strings"

	A "github.com/IBM/fp-go/v2/array"
	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
	Ord "github.com/IBM/fp-go/v2/ord"
	P "github.com/IBM/fp-go/v2/predicate"
)

var stringOrd = Ord.FromCompare(cmp.Compare[string])

// ChooseCluster returns configured when set, the sole discovered name, or an error.
func ChooseCluster(configured string, discovered []string) (string, error) {
	return E.Uneitherize0(func() E.Either[error, string] {
		return F.Pipe2(
			configured,
			O.FromPredicate(P.IsNonZero[string]()),
			O.Fold(
				func() E.Either[error, string] {
					return F.Pipe2(discovered, distinctNonEmpty, chooseDiscovered)
				},
				func(cfg string) E.Either[error, string] { return E.Of[error](cfg) },
			),
		)
	})()
}

func chooseDiscovered(names []string) E.Either[error, string] {
	return A.MatchLeft(
		func() E.Either[error, string] { return E.Left[string](ErrEmptyCluster) },
		func(head string, tail []string) E.Either[error, string] {
			return F.Pipe1(
				tail,
				A.Match(
					func() E.Either[error, string] { return E.Of[error](head) },
					func([]string) E.Either[error, string] { return E.Left[string](ErrAmbiguousCluster) },
				),
			)
		},
	)(names)
}

func distinctNonEmpty(names []string) []string {
	return F.Pipe1(
		names,
		F.Flow4(
			A.Map(strings.TrimSpace),
			A.Filter(P.IsNonZero[string]()),
			strictUniqStrings,
			A.Sort(stringOrd),
		),
	)
}

func strictUniqStrings(names []string) []string {
	return A.StrictUniq(names)
}
