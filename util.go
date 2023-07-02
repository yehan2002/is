package is

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/yehan2002/is/v2/internal"
)

// internal errors used for testing this package
var (
	errNilSuite        = errors.New("test suite is nil")
	errMethodSignature = errors.New("invalid method signature for Setup/Teardown")
	errReceiver        = errors.New("got value not pointer to value")

	errCalledFail    = errors.New("Fail() was called")
	errErrorNotMatch = errors.New("error did not match")
	errFuncNoPanic   = errors.New("function did not panic")
	errNotEqual      = errors.New("values are not equal")
	errCondition     = errors.New("condition was not true")
)

// runT runs the given test function using [*testing.T].
// If the package is being tested, [internal.Test] is used instead.
func runT(t internal.T, opts *options, name string, parallel bool, fn func(Is)) {
	if testingT, ok := t.(*testing.T); ok {
		testingT.Run(name, func(t *testing.T) {
			t.Helper()
			if parallel {
				t.Parallel()
			}

			fn(newIs(t, opts))
		})
	} else if internalT, ok := t.(*internal.Test); ok {
		internalT.Run(name, parallel, func(t *internal.Test) { fn(newIs(t, opts)) })
	}

}

func cmpValue(v1, v2 interface{}, options *options) string {
	return cmp.Diff(v1, v2, options.CmpOpts()...)
}
