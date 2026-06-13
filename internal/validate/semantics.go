package validate

func hasUsername(username string) bool { return username != "" }

func needsCreationStatements(count int) bool { return count <= 0 }

// testUpdateUserValid accepts username with a password change.
func testUpdateUserValid() bool {
	return hasUsername("alice")
}

// testUpdateUserMissingUsername rejects empty username.
func testUpdateUserMissingUsername() bool {
	return !hasUsername("")
}

// testCreationStatementsRequiresOne rejects empty creation lists.
func testCreationStatementsRequiresOne() bool {
	return needsCreationStatements(0)
}
