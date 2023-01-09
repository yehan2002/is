package is

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/yehan2002/is/v2/internal"
)

// Is is provides helpers for writing tests.
type Is func(cond bool, msg string, i ...interface{})

// Equal checks if the given values are equal.
func (is Is) Equal(value, expected interface{}, format string, i ...interface{}) {
	if !reflect.DeepEqual(value, expected) {
		if diff := cmpValue(value, expected); len(diff) != 0 {
			is.t().Helper()
			is.fail(errNotEqual, "Values are not equal:\n"+diff, format, i...)
		}
	}
}

// Fail immediately fails the test.
// Calling this function is the equivalent of calling is.T().Fatalf.
func (is Is) Fail(format string, args ...interface{}) {
	is.t().Helper()
	is.fail(errCalledFail, "", format, args...)
}

// Err checks if any error in err's chain matches target.
// If no errors match target, the test fails.
func (is Is) Err(err, target error, format string, args ...interface{}) {
	if !errors.Is(err, target) {
		is.t().Helper()
		is.fail(errErrorNotMatch, fmt.Sprintf("Error `%s` is not `%s`", err, target), format, args...)
	}
}

// Panic checks if calling the given function causes a panic.
// If the given function does not panic the test fails.
func (is Is) Panic(fn func(), format string, i ...interface{}) {

	var recovered bool

	func() {
		defer func() {
			if r := recover(); r != nil {
				recovered = true
			}
		}()
		fn()
		return
	}()

	if !recovered {
		is.t().Helper()
		is.fail(errFuncNoPanic, "Function did not panic", format, i...)
	}
}

// Log logs the given message.
// This is the equivalent of calling is.T().Log(msg).
// This function can be called from multiple goroutines concurrently.
func (is Is) Log(msg string, i ...interface{}) {
	t := is.t()
	t.Helper()
	t.Logf(msg, i...)
}

// Run runs the given sub test.
// This runs testFn in a separate goroutine and blocks until f returns or calls is.T().Parallel to become a
// parallel test.
func (is Is) Run(name string, testFn func(Is)) { runT(is.t(), name, false, testFn) }

// RunP runs the given test in parallel with the current test.
func (is Is) RunP(name string, testFn func(Is)) { runT(is.t(), name, true, testFn) }

// t gets the underlying internal.T for this test.
func (is Is) t() (t internal.T) {
	// This is a hack to get the internal.T value from `is`.
	// Calling `is(false, "", *internal.T)` sets the the value to the given ptr.
	// This is done by calling `setT`.
	is(false, "", &t)
	return
}

// T gets the underlying *testing.T for this test.
func (is Is) T() *testing.T {
	// this will panic when testing this package because T will be *internal.Test.
	return is.t().(*testing.T)
}

// fail fails the test.
// Calling this function will cause the test to stop executing.
// reason is the reason the test failed. format and i are user provided information about why the
// test failed. The error value passed to this function is only used when testing this package.
func (is Is) fail(err error, reason string, format string, i ...interface{}) {
	t := is.t()
	t.Helper()

	// set the error. This value is used by tests to check if the test failed for the correct reason.
	if internal, ok := t.(*internal.Test); ok {
		internal.SetError(err)
	}

	t.Errorf(format, i...)
	if reason != "" {
		t.Error(reason)
	}

	t.FailNow()
}

func setT(t internal.T, msg string, i []interface{}) (ok bool) {
	if msg == "" && len(i) == 1 {
		var dst *internal.T
		if dst, ok = i[0].(*internal.T); ok {
			*dst = t
		}
		return ok
	}
	return
}

// New creates a new test
func New(t *testing.T) Is { return newIs(t) }

func newIs(t internal.T) Is {
	return func(cond bool, msg string, i ...interface{}) {
		t.Helper()
		if !cond {
			if ok := setT(t, msg, i); ok { // see comment in is.T()
				return
			}

			// set the error. This value is used by tests to check if the test failed for the correct reason.
			if internal, ok := t.(*internal.Test); ok {
				internal.SetError(errCondition)
			}

			t.Fatalf(msg, i...)
		}
	}
}
