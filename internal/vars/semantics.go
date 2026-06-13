package vars

// testNewUserHasRequiredKeys reports ForNewUser supplies every required placeholder.
func testNewUserHasRequiredKeys() bool {
	v := ForNewUser("u", "p", "e", "c")
	return HasRequiredKeys(OpNewUser, v)
}

// testAllOpsHaveRequiredKeys reports each For* builder matches its operation keys.
func testAllOpsHaveRequiredKeys() bool {
	if !HasRequiredKeys(OpNewUser, ForNewUser("u", "p", "e", "c")) {
		return false
	}
	if !HasRequiredKeys(OpUpdatePassword, ForUpdatePassword("u", "p", "c")) {
		return false
	}
	if !HasRequiredKeys(OpUpdateExpiration, ForUpdateExpiration("u", "e", "c")) {
		return false
	}
	return HasRequiredKeys(OpDeleteUser, ForDeleteUser("u"))
}
