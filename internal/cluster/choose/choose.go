package choose

import (
	"sort"
	"strings"
)

// ChooseCluster returns configured when set, the sole discovered name, or an error.
func ChooseCluster(configured string, discovered []string) (string, error) {
	if strings.TrimSpace(configured) != "" {
		return configured, nil
	}
	names := distinctNonEmpty(discovered)
	switch len(names) {
	case 0:
		return "", ErrEmptyCluster
	case 1:
		return names[0], nil
	default:
		return "", ErrAmbiguousCluster
	}
}

func distinctNonEmpty(names []string) []string {
	seen := make(map[string]struct{})
	var out []string
	for _, name := range names {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}
