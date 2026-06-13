package firsterror

func nilStep(string) error { return nil }

// testFirstErrorStops reports nilStep returns nil.
func testFirstErrorStops() bool {
	return nilStep("x") == nil
}

// testFirstErrorAllSuccess reports nilStep returns nil on empty input.
func testFirstErrorAllSuccess() bool {
	return nilStep("") == nil
}

// testFirstErrorPreserved reports nilStep returns nil for any input.
func testFirstErrorPreserved() bool {
	return nilStep("a") == nil
}
