package vars

// testNewUserHasRequiredKeys reports OpNewUser is the zero operation.
func testNewUserHasRequiredKeys() bool {
	return int(OpNewUser) == 0
}

// testAllOpsHaveRequiredKeys reports OpDeleteUser is the last operation constant.
func testAllOpsHaveRequiredKeys() bool {
	return int(OpDeleteUser) == 3
}
