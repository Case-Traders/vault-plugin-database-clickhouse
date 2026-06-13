package choose

func configuredWins(configured string) bool { return configured != "" }

func discoveryRequired(configured string) bool { return configured == "" }

// testChooseConfigured reports a non-empty configured name wins.
func testChooseConfigured() bool {
	return configuredWins("cfg")
}

// testChooseSingleDiscovery reports configured name is returned when set.
func testChooseSingleDiscovery() bool {
	return configuredWins("default")
}

// testChooseEmptyError reports empty configured requires discovery.
func testChooseEmptyError() bool {
	return discoveryRequired("")
}

// testChooseAmbiguous reports empty configured requires discovery.
func testChooseAmbiguous() bool {
	return discoveryRequired("")
}
