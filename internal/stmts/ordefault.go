package stmts

// StatementsOrDefault returns provided when non-empty, otherwise a single fallback statement.
func StatementsOrDefault(provided []string, fallback string) []string {
	if len(provided) > 0 {
		return provided
	}
	return []string{fallback}
}
