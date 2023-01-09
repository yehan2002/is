package is

import (
	"errors"
	"testing"

	"github.com/yehan2002/is/v2/internal"
)

type testTest struct{ testBFails bool }

func (t *testTest) TestA(Is)         {}
func (t *testTest) TestZ(Is)         {}
func (t *testTest) TestInvalidTest() {}
func (t testTest) TestX(Is)          {}
func (t *testTest) TestB(is Is) {
	if t.testBFails {
		is.Fail("failed")
	}
}

type testSetupTeardown struct {
	testTest
	setupCalled bool
}

func (t *testSetupTeardown) Teardown() {}
func (t *testSetupTeardown) Setup()    { t.setupCalled = true }

type testSetupTeardownIncorrect struct{}

func (t *testSetupTeardownIncorrect) Setup() error { return nil }

func TestSuiteNil(t *testing.T) {
	result := internal.Run(func(t internal.T) { makeSuite(t, nil, false) })
	if !result.Failed || !errors.Is(result.TestError, errNilSuite) {
		t.Fatalf("makeSuite allowed a nil test suite: %s", result.FailMessage)
	}
}

func TestSuite(t *testing.T) {
	var suite *testSuite
	var testSuite = &testTest{}

	result := internal.Run(func(t internal.T) { suite = makeSuite(t, testSuite, false) })
	if result.Failed {
		t.Fatal("makeSuite failed for a valid suite")
	}

	tests := []*test{
		{testSuite.TestA, "TestA"},
		{testSuite.TestB, "TestB"},
		{testSuite.TestX, "TestX"},
		{testSuite.TestZ, "TestZ"},
	}

	if len(suite.tests) != len(tests) {
		t.Fatal("Incorrect number of tests")
	}

	for i := range tests {
		if tests[i].Name != suite.tests[i].Name {
			t.Fatal("Incorrect function order")
		}
	}

	result = internal.Run(func(t internal.T) { suite.Run(t) })
	if result.Failed {
		t.Fatal("suite failed")
	}

	expectedRunOrder := []string{"TestA", "TestB", "TestX", "TestZ"}
	for i := range expectedRunOrder {
		if result.RunTests[i].Name != expectedRunOrder[i] {
			t.Fatalf("Expected test %s to be run as the %d test", expectedRunOrder[i], i)
		}
	}
}

func TestSuitFail(t *testing.T) {
	result := internal.Run(func(t internal.T) { makeSuite(t, &testTest{testBFails: true}, false).Run(t) })
	if !result.Failed {
		t.Fatalf("Test should fail")
	}
	if result.FailMessage[0] != "failed" {
		t.Fatalf("Incorrect test message")
	}
}

func TestSuiteSetupTeardown(t *testing.T) {
	var suite *testSuite
	var testSuite = &testSetupTeardown{}

	result := internal.Run(func(t internal.T) { suite = makeSuite(t, testSuite, false); suite.Run(t) })
	if result.Failed {
		t.Fatal("failed to run suite")
	}

	if !testSuite.setupCalled {
		t.Fatal("Setup was not called")
	}

	if len(result.CleanupFuncs) != 1 {
		t.Fatal("Cleanup function was not added")
	}

	result = internal.Run(func(t internal.T) { suite = makeSuite(t, &testSetupTeardownIncorrect{}, false); suite.Run(t) })
	if !result.Failed || !errors.Is(result.TestError, errMethodSignature) {
		t.Fatalf("allowed test suite with invalid setup function: %s", result.FailMessage)
	}
}

func TestSuiteReceiver(t *testing.T) {
	result := internal.Run(func(t internal.T) { makeSuite(t, testTest{}, false) })
	if !result.Failed || !errors.Is(result.TestError, errReceiver) {
		t.Fatalf("allowed suite with pointer receivers to be created from struct value")
	}
}
