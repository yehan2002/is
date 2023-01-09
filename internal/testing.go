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
	Fatalf(format string, args ...interface{})
	Logf(format string, args ...interface{})

	Error(v ...interface{})
	Errorf(format string, args ...interface{})
	FailNow()
}

var _ T = (*testing.T)(nil)

// Test an implementation of [T] used to test the [is] package
type Test struct {
	CleanupFuncs []func()
	RunTests     []TestFn

	Failed      bool
	FailMessage []string
	TestError   error
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
func (t *Test) Run(name string, parallel bool, f func(t *Test)) bool {
	t.RunTests = append(t.RunTests, TestFn{Name: name, Parallel: parallel, F: f})
	f(t)
	return true
}

// TestFn a test function passed to [Test.Run]
type TestFn struct {
	Name     string
	Parallel bool
	F        func(t *Test)
}

var errFatal = errors.New("fatal")

// Fatalf fails the test and panics
func (t *Test) Fatalf(format string, args ...interface{}) {
	t.Failed = true
	t.FailMessage = append(t.FailMessage, fmt.Sprintf(format, args...))
	panic(errFatal)
}

// Error appends an error
func (t *Test) Error(v ...interface{}) {
	t.Failed = true
	t.FailMessage = append(t.FailMessage, fmt.Sprint(v...))
}

// Errorf appends an error
func (t *Test) Errorf(format string, args ...interface{}) {
	t.Failed = true
	t.FailMessage = append(t.FailMessage, fmt.Sprintf(format, args...))
}

// FailNow panic with errFatal
func (t *Test) FailNow() {
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
	t.TestError = err
}
