package stmts

// testStatementsOrDefaultFallback uses fallback when provided is empty.
func testStatementsOrDefaultFallback() bool {
	got := StatementsOrDefault(nil, "ALTER USER")
	return len(got) == 1 && got[0] == "ALTER USER"
}

// testStatementsOrDefaultProvided keeps non-empty provided statements.
func testStatementsOrDefaultProvided() bool {
	in := []string{"CREATE USER", "GRANT"}
	got := StatementsOrDefault(in, "ALTER USER")
	if len(got) != len(in) {
		return false
	}
	for i := range in {
		if got[i] != in[i] {
			return false
		}
	}
	return true
}
