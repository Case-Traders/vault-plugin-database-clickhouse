package stmt

import "strings"

// NormalizeCommands splits each command on semicolons, trims whitespace, and drops empty parts.
func NormalizeCommands(commands []string) []string {
	var out []string
	for _, cmd := range commands {
		out = append(out, normalizeStatement(cmd)...)
	}
	return out
}

func normalizeStatement(stmt string) []string {
	var out []string
	for _, part := range splitSemicolon(stmt) {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}

func splitSemicolon(s string) []string {
	return strings.Split(s, ";")
}
