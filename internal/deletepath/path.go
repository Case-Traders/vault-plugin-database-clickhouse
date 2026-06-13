package deletepath

// UseCustomRevocation reports whether DeleteUser should run role revocation statements.
func UseCustomRevocation(commands []string) bool {
	return len(commands) > 0
}
