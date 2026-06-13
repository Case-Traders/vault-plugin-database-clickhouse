package vars

import (
	"testing"

	"pgregory.net/rapid"
)

func TestHasRequiredKeys_newUser(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		v := ForNewUser(
			rapid.String().Draw(t, "user"),
			rapid.String().Draw(t, "pass"),
			rapid.String().Draw(t, "exp"),
			rapid.String().Draw(t, "cluster"),
		)
		if !HasRequiredKeys(OpNewUser, v) {
			t.Fatalf("missing keys in %v", v.Keys())
		}
	})
}

func TestHasRequiredKeys_allOperations(t *testing.T) {
	ops := []struct {
		op Operation
		v  TemplateVars
	}{
		{OpNewUser, ForNewUser("u", "p", "e", "c")},
		{OpUpdatePassword, ForUpdatePassword("u", "p", "c")},
		{OpUpdateExpiration, ForUpdateExpiration("u", "e", "c")},
		{OpDeleteUser, ForDeleteUser("u")},
	}
	for _, tc := range ops {
		if !HasRequiredKeys(tc.op, tc.v) {
			t.Fatalf("op %d missing keys %v", tc.op, RequiredKeys(tc.op))
		}
	}
}
