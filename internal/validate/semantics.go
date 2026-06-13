package validate

// testUpdateUserValid accepts username with a password change.
func testUpdateUserValid() bool {
	return UpdateUser("alice", true, false) == nil
}

// testUpdateUserMissingUsername rejects empty username.
func testUpdateUserMissingUsername() bool {
	return UpdateUser("", true, false) != nil
}

// testCreationStatementsRequiresOne rejects empty creation lists.
func testCreationStatementsRequiresOne() bool {
	return CreationStatements(0) != nil
}
