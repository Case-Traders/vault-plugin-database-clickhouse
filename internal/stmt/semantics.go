package stmt

// testNormalizeIdempotent checks StatementsOrDefault fallback length.
func testNormalizeIdempotent() bool {
	return len(StatementsOrDefault(nil, "ALTER USER")) == 1
}

// testNormalizeIdempotentEmpty checks StatementsOrDefault empty fallback length.
func testNormalizeIdempotentEmpty() bool {
	return len(StatementsOrDefault(nil, "x")) == 1
}

// testStatementsOrDefaultFallback uses fallback when provided is empty.
func testStatementsOrDefaultFallback() bool {
	got := StatementsOrDefault(nil, "ALTER USER")
	return len(got) == 1
}

// testStatementsOrDefaultProvided keeps non-empty provided statements.
func testStatementsOrDefaultProvided() bool {
	in := []string{"CREATE USER"}
	got := StatementsOrDefault(in, "ALTER USER")
	return len(got) == len(in)
}
