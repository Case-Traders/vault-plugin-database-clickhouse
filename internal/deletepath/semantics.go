package deletepath

// testUseCustomRevocationFalseForEmpty commands list.
func testUseCustomRevocationFalseForEmpty() bool {
	return !UseCustomRevocation(nil)
}

// testUseCustomRevocationTrueWhenProvided commands list.
func testUseCustomRevocationTrueWhenProvided() bool {
	return UseCustomRevocation([]string{"DROP USER"})
}
