package txexec

import "errors"

var errStop = errors.New("stop")

// testFirstErrorStops reports execution stops at the first failing step.
func testFirstErrorStops() bool {
	calls := 0
	err := FirstError([]string{"a", "b", "c"}, func(s string) error {
		calls++
		if s == "b" {
			return errStop
		}
		return nil
	})
	return err == errStop && calls == 2
}

// testFirstErrorAllSuccess reports nil when every step succeeds.
func testFirstErrorAllSuccess() bool {
	err := FirstError([]string{"a", "b", "c"}, func(string) error { return nil })
	return err == nil
}

// testFirstErrorPreserved reports the first error is returned unchanged.
func testFirstErrorPreserved() bool {
	want := errors.New("boom")
	err := FirstError([]string{"1", "2", "3"}, func(s string) error {
		if s == "1" {
			return want
		}
		return errors.New("later")
	})
	return err == want
}
