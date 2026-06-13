package stmt

import "slices"

// testNormalizeIdempotent checks NormalizeCommands is idempotent (Goose semantics test).
func testNormalizeIdempotent() bool {
	in := []string{"CREATE USER; GRANT x;", ";;  DROP;"}
	n1 := NormalizeCommands(in)
	n2 := NormalizeCommands(n1)
	return slices.Equal(n1, n2)
}

// testNormalizeIdempotentEmpty checks idempotence on an empty input slice.
func testNormalizeIdempotentEmpty() bool {
	once := NormalizeCommands(nil)
	return slices.Equal(NormalizeCommands(once), once)
}
