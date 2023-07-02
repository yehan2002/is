package is

import (
	"reflect"
	"runtime/debug"
	"strings"
	"testing"

	"github.com/yehan2002/is/v2/internal"
)

// Suite runs the given test suite.
// All test functions must have the prefixed with Test and must take [Is] as the first and only argument
// and should not have a return value.
// Test functions are run sequentially in lexicographic order. To run tests in parallel use [SuiteP].
//
//	// Example test
//	func (s *SuiteName) TestName(Is){}
//
// A test suite can define a Setup and Teardown function to setup the test suite and cleanup after the
// suite has completed. Both setup and teardown functions must have no arguments and return values.
// The Teardown function is always called after running all tests even if the tests fail.
//
//	func (s *suiteName) Setup(){ /* setup the set suite here */ }
//	func (s *suiteName) Teardown(){ /* clean up after the test suite has completed. */ }
func Suite(t *testing.T, suite interface{}, opts ...Option) {
	t.Helper()
	makeSuite(t, suite, false, opts).Run(t)
}

// SuiteP like [Suite] but calls all test function in parallel.
func SuiteP(t *testing.T, suite interface{}, opts ...Option) {
	t.Helper()
	makeSuite(t, suite, true, opts).Run(t)
}

type testSuite struct {
	name     string
	parallel bool

	setupFunc    func()
	teardownFunc func()

	options *options

	tests []*test
}

type test struct {
	Func func(Is)
	Name string
}

func (s *testSuite) Run(t internal.T) {
	t.Helper()

	// skip suite if it has no tests
	if len(s.tests) == 0 {
		t.Logf("is.Suite: skipped suite '%s' with no tests", s.name)
		return
	}

	s.setupFunc()
	t.Cleanup(s.teardownFunc)

	for i := range s.tests {
		test := s.tests[i]
		runT(t, s.options, test.Name, s.parallel, test.Func)
	}
}

func makeSuite(t internal.T, s interface{}, parallel bool, opts []Option) (testS *testSuite) {
	// calledFatal indicates that t.Fatal was called.
	// This is used to differentiate between t.Fatal calling runtime.Goexit and a panic in the code bellow.
	var calledFatal bool

	defer func() {
		// t.Fatal was called. don't try to recover.
		if calledFatal {
			return
		}

		if r := recover(); r != nil {
			t.Fatalf("is.Suite: Internal error: %s\n%s", r, debug.Stack())
		}
	}()

	fatal := func(err error, f string, args ...interface{}) {
		t.Helper()
		calledFatal = true

		// this package is being tested, set the error to be verified by the test.
		if internal, ok := t.(*internal.Test); ok {
			internal.SetError(err)
		}

		t.Fatalf(f, args...)
	}

	options := newOptions(opts)

	t.Helper()

	suite := reflect.ValueOf(s)
	if isNil(suite) {
		fatal(errNilSuite, "is.Suite: test suite is nil.")
	}

	suiteType := suite.Type()
	testS = &testSuite{name: suite.Type().Name(), parallel: parallel, options: options}

	// check if the caller passed a value instead of a pointer to a value by accident.
	// If any methods on the type have a pointer receiver, they cannot be called because `suite` is not
	// addressable. Creating a addressable copy of suite is not safe because caller may expect the value
	// of suite to be modified by running the test.
	if suitePtr := reflect.PointerTo(suite.Type()); suitePtr.NumMethod() != 0 {
		for i := 0; i < suitePtr.NumMethod(); i++ {
			method := suitePtr.Method(i)
			methodType := method.Type

			if methodType.NumIn() == 0 {
				// this should never happen because a method will always have at least one argument (the receiver).
				continue
			}

			// check if the method has a pointer receiver.
			if methodType.In(0) == suitePtr {
				if n := method.Name; n == "Setup" || n == "Teardown" || strings.HasPrefix(n, "Test") {
					fatal(errReceiver, "is.Suite: Method %s has a pointer receiver but Suite was given a %s not *%s.", n, testS.name, testS.name)
				}
			}
		}
	}

	// get setup and teardown functions
	testS.setupFunc = getMethod(fatal, suite, "Setup")
	testS.teardownFunc = getMethod(fatal, suite, "Teardown")

	// get all tests defined by the suite.
	for i := 0; i < suite.NumMethod(); i++ {
		methodValue := suite.Method(i)
		name := suiteType.Method(i).Name

		// ignore unexported methods
		if !strings.HasPrefix(name, "Test") || !methodValue.CanInterface() {
			continue
		}

		testFunc, ok := methodValue.Interface().(func(Is))
		if !ok {
			t.Logf("is.Suite: Skipping test function '%s' with incorrect method signature. Should be func(Is) ", name)
			continue
		}

		testS.tests = append(testS.tests, &test{Name: name, Func: testFunc})
	}

	return
}

// getMethod gets the method of v that has the given name.
// If the method does not exist, an no-op function is returned instead.
// This function calls fatal if the method exists but isn't the same type as F.
func getMethod(fatal func(err error, f string, a ...interface{}), v reflect.Value, name string) (F func()) {
	method := v.MethodByName(name)
	if method.IsValid() {
		Func, ok := method.Interface().(func())
		if !ok {
			fatal(errMethodSignature, "is.Suite: %s method should have no arguments and return values", name)
		}
		return Func
	}

	return func() {}
}

func isNil(v reflect.Value) (isNil bool) {
	if !v.IsValid() {
		return true
	}

	switch v.Kind() {
	case reflect.Func, reflect.Map, reflect.Pointer, reflect.Interface, reflect.Slice, reflect.UnsafePointer, reflect.Chan:
		return v.IsNil()
	default:
		return false
	}
}
