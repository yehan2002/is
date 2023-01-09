package is

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/yehan2002/is/v2/internal"
)

// internal errors used for tests
var (
	errNilSuite        = errors.New("test suite is nil")
	errMethodSignature = errors.New("invalid method signature for Setup/Teardown")
	errReceiver        = errors.New("got value not pointer to value")

	errCalledFail    = errors.New("Fail() was called")
	errErrorNotMatch = errors.New("error did not match")
	errFuncNoPanic   = errors.New("function did not panic")
	errNotEqual      = errors.New("values are not equal")
	errCondition     = errors.New("condition was not true")
)

func runT(t internal.T, name string, parallel bool, fn func(Is)) {
	if testingT, ok := t.(*testing.T); ok {
		testingT.Run(name, func(t *testing.T) {
			t.Helper()
			if parallel {
				t.Parallel()
			}

			fn(newIs(t))
		})
	} else if internalT, ok := t.(*internal.Test); ok {
		internalT.Run(name, parallel, func(t *internal.Test) { fn(newIs(t)) })
	}

}

func cmpValue(v1, v2 interface{}) string {
	return cmp.Diff(v1, v2, cmp.FilterPath(func(p cmp.Path) bool {
		sf, ok := p.Index(-1).(cmp.StructField)
		if !ok {
			return false
		}

		field := p.Index(-2).Type().Field(sf.Index())
		isExported := field.PkgPath == ""

		return !isExported || field.Tag != "" &&
			(field.Tag.Get("deep") == "-" || field.Tag.Get("cmp") == "-")
	}, cmp.Ignore()))
}
