package is

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"
)

//Suite runs a test suite.
//If the suite contains a method named `Setup` it is called before any tests are run.
//Tests must start with `Test` and take `is.IS` as the first arg.
//Finally after all the tests are called `Teardown` is called
func Suite(t *testing.T, v interface{}) {
	if v == nil {
		panic("The provided suite is nil")
	}

	is := &baseTest{t: t}
	isr := reflect.ValueOf(is)
	s := &testSuite{is: is, isr: isr, passed: true, testChan: make(chan bool, 1), t: t}

	if reflect.TypeOf(v).Kind() == reflect.Ptr {
		s.name = reflect.TypeOf(v).Elem().Name()
	} else {
		s.name = reflect.TypeOf(v).Name()
	}
	printf(messages.start, true, s.name)

	is.fail = s.fail
	var suite, suitePtr reflect.Value

	if reflect.TypeOf(v).Kind() == reflect.Ptr {
		suitePtr, suite = reflect.ValueOf(v), reflect.ValueOf(v).Elem()
	} else {
		suitePtr, suite = reflect.New(reflect.TypeOf(v)), reflect.ValueOf(v)
		suitePtr.Elem().Set(reflect.ValueOf(v))
	}

	s.setupTests(suitePtr)
	s.callTests(suite)
	s.callTests(suitePtr)
	s.teardownTests(suitePtr)

	printf(messages.passSuite, true, s.name)
}

type testSuite struct {
	is       IS
	isr      reflect.Value
	passed   bool
	testChan chan bool
	t        *testing.T
	name     string
	color    bool
}

func (t *testSuite) callTests(value reflect.Value) {
	for i := 0; i < value.NumMethod(); i++ {
		method := value.Type().Method(i)
		if strings.HasPrefix(method.Name, "Test") {
			if method.Type.NumIn() == 2 && method.Type.In(1) == isType {
				t.callTest(method, value)
			} else {
				fmt.Printf("Skipping method \"%s\" with invalid method signature", method.Name)
			}
		}
	}
}

func (t *testSuite) callTest(method reflect.Method, value reflect.Value) {
	now := time.Now()
	if testing.Verbose() {
		printf(messages.run, true, method.Name)
	}
	stdout := captureStdout()
	if t.callTestFunc(method, value) && t.passed {
		buf := stdout()
		if testing.Verbose() { // ignore stdout since the test passed
			fmt.Println(buf)
		}

		printf(messages.passTest, true, time.Now().Sub(now).Seconds())
	} else {
		buf := stdout()
		printf(messages.failTest, true, time.Now().Sub(now).Seconds())
		fmt.Println(buf)
		printf(messages.failSuite, true, t.name)
		t.t.FailNow()
	}

}

func (t *testSuite) callTestFunc(method reflect.Method, value reflect.Value) (ok bool) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				t.is.Fail(fmt.Sprintf("panic: %s", r))
				t.testChan <- false
			}
		}()
		method.Func.Call([]reflect.Value{value, t.isr})
		t.testChan <- true
	}()
	pass, ok := <-t.testChan
	return pass && ok
}

func (t *testSuite) setupTests(value reflect.Value) {
	if method, ok := value.Type().MethodByName("Setup"); ok {
		if method.Type.NumIn() != 1 {
			panic("Invalid method signature for Setup")
		}
		method.Func.Call([]reflect.Value{value})
	}
}

func (t *testSuite) teardownTests(value reflect.Value) {
	if method, ok := value.Type().MethodByName("Teardown"); ok {
		if method.Type.NumIn() != 1 {
			panic("Invalid method signature for Teardown")
		}
		method.Func.Call([]reflect.Value{value})
	}
}
