package is

import (
	"reflect"
	"strings"
	"testing"
)

type suiteTest struct {
	fn   func(Is)
	name string
}

// Suite runs the given test suite.
// If the suite contains a method named `Setup` it is called before any tests are run.
// Tests must start with `Test` and take `is.IS` as the first arg.
// Finally after all the tests are run the `Teardown` method is called.
func Suite(t *testing.T, suite interface{}) {
	t.Helper()
	runSuite(t, suite, false)
}

// SuiteP like `Suite` but calls all test function in parallel.
func SuiteP(t *testing.T, suite interface{}) {
	t.Helper()
	runSuite(t, suite, true)
}

func runSuite(t *testing.T, s interface{}, parallel bool) {
	t.Helper()
	suite := reflect.ValueOf(s)
	if k := suite.Kind(); k != reflect.Ptr && k != reflect.Interface {
		t.Fatal("is.Suite: suite must be a ptr or interface")
	}

	suiteType := suite.Type()

	var tests []*suiteTest
	var setup = func() {}
	var teardown = func() {}
	var ok bool

	for i := 0; i < suite.NumMethod(); i++ {
		methodType := suiteType.Method(i)
		method := suite.Method(i)
		if !isExported(methodType) {
			continue
		}
		if name := methodType.Name; strings.HasPrefix(name, "Test") {
			if fn, ok := method.Interface().(func(Is)); ok {
				tests = append(tests, &suiteTest{name: name, fn: fn})
				continue
			}
			t.Logf("is.Suite: Skipping test function '%s' with incorrect method signature. Should be func(Is) ", name)
			continue
		} else if name == "Setup" {
			if setup, ok = method.Interface().(func()); !ok {
				t.Fatal("is.Suite: Setup function should be have no args and no return values")
			}
		} else if name == "Teardown" {
			if teardown, ok = method.Interface().(func()); !ok {
				t.Fatal("is.Suite: Teardown function should be have no args and no return values")
			}
		}
	}

	if len(tests) == 0 {
		t.Fatalf("is.Suite: skipped suite '%s' with no tests", suiteType.Name())
	}

	setup()
	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			t.Helper()

			if parallel {
				t.Parallel()
			}

			test.fn(New(t))
		})
	}
	t.Cleanup(teardown)
}
