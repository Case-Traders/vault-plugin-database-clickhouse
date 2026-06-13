package cluster

import (
	"errors"
	"strings"
	"testing"

	"pgregory.net/rapid"
)

func nonEmptyString(t *rapid.T, label string) string {
	return rapid.String().Filter(func(s string) bool { return strings.TrimSpace(s) != "" }).Draw(t, label)
}

func TestChooseCluster_configured(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		name := nonEmptyString(t, "name")
		disc := rapid.SliceOf(rapid.String()).Draw(t, "disc")
		got, err := ChooseCluster(name, disc)
		if err != nil || got != name {
			t.Fatalf("configured=%q disc=%v got=%q err=%v", name, disc, got, err)
		}
	})
}

func TestChooseCluster_empty(t *testing.T) {
	_, err := ChooseCluster("", nil)
	if !errors.Is(err, ErrEmptyCluster) {
		t.Fatalf("want ErrEmptyCluster, got %v", err)
	}
}

func TestChooseCluster_single(t *testing.T) {
	got, err := ChooseCluster("", []string{"default"})
	if err != nil || got != "default" {
		t.Fatalf("got=%q err=%v", got, err)
	}
}

func TestChooseCluster_ambiguous(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		c1 := nonEmptyString(t, "c1")
		c1Trim := strings.TrimSpace(c1)
		c2 := rapid.String().Filter(func(s string) bool {
			s = strings.TrimSpace(s)
			return s != "" && s != c1Trim
		}).Draw(t, "c2")
		_, err := ChooseCluster("", []string{c1, c2})
		if !errors.Is(err, ErrAmbiguousCluster) {
			t.Fatalf("want ErrAmbiguousCluster for %q and %q, got %v", c1, c2, err)
		}
	})
}
