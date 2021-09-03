package is

import (
	"fmt"
	"strings"
	"testing"

	"github.com/go-test/deep"
)

// Is is provides helpers for writing tests.
type Is func(cond bool, msg string, i ...interface{})

// Equal checks if the given values are equal
func (is Is) Equal(v1, v2 interface{}, msg string, i ...interface{}) {
	if dif := deep.Equal(v1, v2); len(dif) != 0 {
		is.T().Helper()
		is(false, fmt.Sprintf("%s\nValues are not equal:\n\t%s", fmt.Sprintf(msg, i...), strings.Join(dif, "\n\t")))
	}
}

// Fail immediately fails the test.
func (is Is) Fail(msg string, i ...interface{}) {
	is.T().Helper()
	is(false, msg, i...)
}

// Panic checks if calling the given function causes a panic.
// If the given function does not panic the test fails.
func (is Is) Panic(panicable func(), msg string, i ...interface{}) {
	if !callPanic(panicable) {
		is.T().Helper()
		is(false, fmt.Sprintf("%s\nFunction did not panic", fmt.Sprintf(msg, i...)))
	}
}

// Log logs the given message.
// This is the equivalent of calling is.T().Log(msg).
// This function can be called from multiple goroutines concurrently.
func (is Is) Log(msg string, i ...interface{}) {
	t := is.T()
	t.Helper()
	t.Logf(msg, i...)
}

// Run runs the given test.
func (is Is) Run(name string, f func(Is)) {
	is.T().Run(name, func(t *testing.T) { f(New(t)) })
}

// RunP runs the given test in parallel with the current test.
func (is Is) RunP(name string, f func(Is)) {
	is.T().Run(name, func(t *testing.T) { t.Parallel(); f(New(t)) })
}

// T gets the underlying *testing.T for this test.
func (is Is) T() (t *testing.T) {
	// This is a ugly hack to get the testing.T value from `is`.
	// Calling `is(false, "", internalIsCall, **testing.T)` sets the the value to the given ptr.
	// This is done by `setT`
	is(false, "", internalIsCall, &t)
	return
}

var internalIsCall = new(uint16)

func setT(t *testing.T, msg string, i []interface{}) (ok bool) {
	if msg == "" && len(i) == 2 {
		if i[0] == internalIsCall {
			var dst **testing.T
			if dst, ok = i[1].(**testing.T); ok {
				*dst = t
			}
			return ok
		}
	}
	return
}

func callPanic(f func()) (paniced bool) {
	defer func() {
		if r := recover(); r != nil {
			paniced = true
		}
	}()
	f()
	return
}

// New creates a new test
func New(t *testing.T) Is {
	return func(cond bool, msg string, i ...interface{}) {
		t.Helper()
		if !cond {
			if ok := setT(t, msg, i); ok { // see comment in is.T()
				return
			}
			t.Errorf(msg, i...)
			t.FailNow()
		}
	}
}
