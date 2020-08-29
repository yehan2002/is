package is

import (
	"fmt"
	"runtime"
	"testing"
)

type failable func(t *testing.T, test interface{}, comment bool, msg []interface{})

func (t *testSuite) fail(_ *testing.T, test interface{}, comment bool, msg []interface{}) {
	t.passed = false
	if test != nil {
		printf(messages.err2, true, test)
	}
	if msg != nil && len(msg) > 0 {
		printf(messages.err1, true, fmt.Sprint(msg...))
	} else if comment {
		if c, ok := getComment(); ok {
			printf(messages.err1, true, c)
		}
	}

	fmt.Println(getStack(3))
	t.testChan <- false
	runtime.Goexit()
}

func basicFailable(t *testing.T, test interface{}, comment bool, msg []interface{}) {
	if test != nil {
		printf(messages.err2, true, test)
	}

	if msg != nil && len(msg) > 0 {
		printf(messages.err1, true, fmt.Sprint(msg...))
	} else if comment {
		if c, ok := getComment(); ok {
			printf(messages.err1, true, c)
		}
	}

	fmt.Println(getStack(3))
	t.FailNow()
}
