package txexec

import (
	"errors"
	"fmt"
	"testing"

	"pgregory.net/rapid"
)

func TestFirstError_stopsAtFirstFailure(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		n := rapid.IntRange(1, 30).Draw(t, "n")
		failAt := rapid.IntRange(0, n-1).Draw(t, "failAt")
		queries := make([]int, n)
		for i := range queries {
			queries[i] = i
		}
		var ran []int
		err := FirstError(queries, func(q int) error {
			ran = append(ran, q)
			if q == failAt {
				return fmt.Errorf("fail at %d", q)
			}
			return nil
		})
		if err == nil {
			t.Fatal("expected error")
		}
		if len(ran) != failAt+1 {
			t.Fatalf("failAt=%d ran=%d queries %v", failAt, len(ran), ran)
		}
	})
}

func TestFirstError_allSuccess(t *testing.T) {
	err := FirstError([]string{"a", "b", "c"}, func(string) error { return nil })
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFirstError_firstErrorPreserved(t *testing.T) {
	want := errors.New("boom")
	err := FirstError([]int{1, 2, 3}, func(i int) error {
		if i == 1 {
			return want
		}
		return errors.New("later")
	})
	if !errors.Is(err, want) {
		t.Fatalf("got %v want %v", err, want)
	}
}
