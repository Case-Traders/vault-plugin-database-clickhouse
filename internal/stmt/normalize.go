package stmt

import (
	"strings"

	A "github.com/IBM/fp-go/v2/array"
	F "github.com/IBM/fp-go/v2/function"
	P "github.com/IBM/fp-go/v2/predicate"
	"github.com/hashicorp/go-secure-stdlib/strutil"
)

// NormalizeCommands splits each command on semicolons, trims whitespace, and drops empty parts.
func NormalizeCommands(commands []string) []string {
	return F.Pipe1(commands, A.Chain(normalizeStatement))
}

func normalizeStatement(stmt string) []string {
	return F.Pipe2(
		strutil.ParseArbitraryStringSlice(stmt, ";"),
		A.Map(strings.TrimSpace),
		A.Filter(P.IsNonZero[string]()),
	)
}
