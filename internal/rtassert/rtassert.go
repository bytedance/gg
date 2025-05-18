// Package rtassert provides runtime assertion.
package rtassert

import (
	"fmt"

	"github.com/bytedance/gg/internal/constraints"
)

func MustNotNeg[T constraints.Number](n T) {
	if n < 0 {
		panic(fmt.Errorf("number must not be negative: %v", n))
	}
}

func ErrMustNil(err error) {
	if err != nil {
		panic(fmt.Errorf("unexpected error: %s", err))
	}
}
