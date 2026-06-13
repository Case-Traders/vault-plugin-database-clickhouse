package choose

// testChooseConfigured reports configured cluster wins over discovery (Goose semantics test).
func testChooseConfigured() bool {
	got, err := ChooseCluster("cfg", []string{"a", "b"})
	return err == nil && got == "cfg"
}

// testChooseSingleDiscovery reports a lone discovered cluster is chosen.
func testChooseSingleDiscovery() bool {
	got, err := ChooseCluster("", []string{"default"})
	return err == nil && got == "default"
}

// testChooseEmptyError reports no clusters yields ErrEmptyCluster.
func testChooseEmptyError() bool {
	_, err := ChooseCluster("", nil)
	return err == ErrEmptyCluster
}

// testChooseAmbiguous reports two distinct clusters yield ErrAmbiguousCluster.
func testChooseAmbiguous() bool {
	_, err := ChooseCluster("", []string{"alpha", "beta"})
	return err == ErrAmbiguousCluster
}
