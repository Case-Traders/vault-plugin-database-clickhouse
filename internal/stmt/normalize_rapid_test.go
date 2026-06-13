package stmt

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"testing"

	"pgregory.net/rapid"
)

func TestNormalize_idempotent(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		n := rapid.IntRange(0, 20).Draw(t, "n")
		commands := make([]string, n)
		for i := range commands {
			commands[i] = rapid.String().Draw(t, "cmd")
		}
		once := NormalizeCommands(commands)
		twice := NormalizeCommands(once)
		if !slices.Equal(once, twice) {
			t.Fatalf("idempotent violated: %v -> %v -> %v", commands, once, twice)
		}
	})
}

func TestNormalize_goldenVectors(t *testing.T) {
	path := goldenPath(t, "stmt_normalize.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read golden vectors: %v", err)
	}
	var cases []struct {
		Input  []string `json:"input"`
		Output []string `json:"output"`
	}
	if err := json.Unmarshal(data, &cases); err != nil {
		t.Fatalf("parse golden vectors: %v", err)
	}
	for i, tc := range cases {
		got := NormalizeCommands(tc.Input)
		if !slices.Equal(got, tc.Output) {
			t.Fatalf("case %d: input=%v got=%v want=%v", i, tc.Input, got, tc.Output)
		}
	}
}

func goldenPath(t *testing.T, name string) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Join(filepath.Dir(file), "..", "..", "proof", "testvectors", name)
}
