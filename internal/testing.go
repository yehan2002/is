package internal

import (
	"errors"
	"fmt"
	"testing"
)

// T is an interface implemented by [testing.T] and Test.
type T interface {
	Helper()
	Parallel()
	Cleanup(f func())
	Run(name string, f func(t *testing.T)) bool
	Fatalf(format string, args ...interface{})
	Logf(format string, args ...interface{})
}

var _ T = (*testing.T)(nil)

// Test an implementation of [T] used to test the [is] package
type Test struct {
	CleanupFuncs []func()
	RunTests     []TestFn

	Failed      bool
	FailMessage string
	Error       error
}

// Helper is a no-op function
func (t *Test) Helper() {}

// Parallel is a no-op function
func (t *Test) Parallel() {}

// Logf is a no-op function
func (t *Test) Logf(format string, args ...interface{}) {}

// Cleanup registers a cleanup function
func (t *Test) Cleanup(f func()) { t.CleanupFuncs = append(t.CleanupFuncs, f) }

// Run adds a test to be run
func (t *Test) Run(name string, f func(t *testing.T)) bool {
	t.RunTests = append(t.RunTests, TestFn{Name: name, F: f})

	return true
}

// TestFn a test function passed to [Test.Run]
type TestFn struct {
	Name string
	F    func(t *testing.T)
}

var errFatal = errors.New("fatal")

// Fatalf fails the test and panics
func (t *Test) Fatalf(format string, args ...interface{}) {
	t.Failed = true
	t.FailMessage = fmt.Sprintf(format, args...)
	panic(errFatal)
}

// Run runs the given test
func Run(f func(T)) (t *Test) {
	t = &Test{}

	defer func() {
		if t.Failed {
			recover()
		}
	}()

	f(t)
	return
}

// SetError sets the error that occurred
func (t *Test) SetError(err error) {
	t.Error = err
}
