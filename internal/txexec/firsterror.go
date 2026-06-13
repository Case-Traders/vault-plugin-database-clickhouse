package txexec

import (
	A "github.com/IBM/fp-go/v2/array"
	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
)

// FirstError runs f on each element in order and returns the first error from f.
func FirstError[T any](xs []T, f func(T) error) error {
	return F.Pipe2(
		xs,
		A.FindFirstMap(func(x T) O.Option[error] {
			if err := f(x); err != nil {
				return O.Some(err)
			}
			return O.None[error]()
		}),
		O.GetOrElse(F.Constant[error](nil)),
	)
}
