package is

import (
	"fmt"
	"runtime"
	"testing"
)

type failable func(t *testing.T, msg interface{}, test interface{}, comment bool)

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

func basicFailable(t *testing.T, msg interface{}, test interface{}, comment bool) {
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
	t.FailNow()
}
