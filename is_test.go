package is

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/yehan2002/is/v2/internal"
)

type deepIgnore struct {
	V int `deep:"-"`
}

type cmpIgnore struct {
	V int `cmp:"-"`
}

// mustFail tests if calling fn causes the test to fail with the given error
func mustFail(t *testing.T, err error, fn func(is Is)) {
	result := internal.Run(func(t internal.T) { fn(newIs(t)) })
	if !result.Failed {
		t.Fatal("Test did not fail")
	}
	if !errors.Is(result.TestError, err) {
		t.Fatalf("Test failed with %s not %s", result.TestError, err)
	}
}

// mustPass tests if calling fn does not cause the test to fail
func mustPass(t *testing.T, fn func(is Is)) {
	result := internal.Run(func(t internal.T) { fn(newIs(t)) })
	if result.Failed {
		t.Fatalf("Test failed: %s", strings.Join(result.FailMessage, "\n"))
	}
}

func TestIs(t *testing.T) {
	mustFail(t, errCondition, func(is Is) { is(false, "this should fail") })
	mustPass(t, func(is Is) { is(true, "this should pass") })

	mustFail(t, errNotEqual, func(is Is) { is.Equal(false, true, "this should fail") })
	mustPass(t, func(is Is) { is.Equal(true, true, "this should pass") })
	mustPass(t, func(is Is) {
		var z, z2 int64
		is.Equal(&z, &z2, "this should pass")
	})
	mustFail(t, errNotEqual, func(is Is) {
		var z, z2 int64
		z2 = 12
		is.Equal(&z, &z2, "this must fail")
	})
	mustPass(t, func(is Is) {
		var z internal.Test
		is.Equal(&z, &z, "this should pass")
	})
	mustFail(t, errNotEqual, func(is Is) {
		var z, z2 internal.Test
		z2.TestError = errFuncNoPanic
		is.Equal(&z, &z2, "this must fail")
	})
	mustFail(t, errNotEqual, func(is Is) {
		var z internal.Test
		var z2 internal.Test
		z2.TestError = errFuncNoPanic
		is.Equal(&z, &z2, "this must fail")
	})
	mustPass(t, func(is Is) {
		var z, z2 deepIgnore
		z.V = 1
		z2.V = 1000
		is.Equal(&z, &z2, "this should pass")
	})
	mustPass(t, func(is Is) {
		var z, z2 cmpIgnore
		z.V = 1
		z2.V = 1000
		is.Equal(&z, &z2, "this should pass")
	})

	mustFail(t, errCalledFail, func(is Is) { is.Fail("this should fail") })

	mustFail(t, errErrorNotMatch, func(is Is) { is.Err(errCalledFail, os.ErrClosed, "this should fail") })
	mustPass(t, func(is Is) { is.Err(os.ErrClosed, os.ErrClosed, "this should pass") })

	mustFail(t, errFuncNoPanic, func(is Is) { is.Panic(func() {}, "this should fail") })
	mustPass(t, func(is Is) { is.Panic(func() { panic("err") }, "this should pass") })

	mustPass(t, func(is Is) {
		is.Log("test")
		is.Run("name", func(i Is) {})
		is.RunP("name", func(i Is) {})
	})
}
