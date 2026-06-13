package firsterror

// FirstError runs f on each string in order and returns the first error from f.
func FirstError(xs []string, f func(string) error) error {
	for _, x := range xs {
		if err := f(x); err != nil {
			return err
		}
	}
	return nil
}
