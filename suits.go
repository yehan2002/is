package is

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"
)

var noColorFlag bool

func init() {
	envNoColor := os.Getenv("NO_COLOR") == "true"
	flag.BoolVar(&noColorFlag, "nocolor", envNoColor, "turns off colors")
}

//isType the reflect.Type of `IS`
var isType = reflect.TypeOf((*IS)(nil)).Elem()

var stdoutMux = &sync.Mutex{}

//Suite runs a test suite
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
	if testing.Verbose() || true {
		printf(messages.run, true, method.Name)
	}
	stdout := captureStdout()
	if t.callTestFunc(method, value) && t.passed {
		if testing.Verbose() || true {
			stdout() // ignore stdout since the test passed
			printf(messages.passTest, true, time.Now().Sub(now).Seconds())
		}
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

func (t *testSuite) fail(_ *testing.T, msg interface{}, test interface{}, comment bool) {
	t.passed = false
	if test != nil {
		printf(messages.err2, true, test)
	}
	if msg != nil {
		printf(messages.err1, true, msg)
	}
	if comment {
		if c, ok := getComment(); ok {
			printf(messages.err1, true, c)
		}
	}

	fmt.Println(getStack(3))
	t.testChan <- false
	runtime.Goexit()
}
