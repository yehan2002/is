package is

import (
	"reflect"
	"testing"
)

// IS is a package for writing tests

var isType = reflect.TypeOf((*IS)(nil)).Elem()

type panicable func()

//New creates a new test helper
func New(t *testing.T) IS {
	return &baseTest{t: t, fail: basicFailable}
}

//IS a test
type IS interface {
	Equal(v1, v2 interface{}) IS
	EqualM(v1, v2 interface{}, msg string) IS

	NotEqual(v1, v2 interface{}) IS
	NotEqualM(v1, v2 interface{}, msg string) IS
	NotNil(v1 interface{}) IS
	NotNilM(v1 interface{}, msg string) IS
	Nil(v1 interface{}) IS
	NilM(v1 interface{}, msg string) IS

	MustPanic()
	MustPanicCall(panicable)
	MustPanicCallReflect(funct interface{}, args ...interface{})

	Err(v1 interface{}) IS
	True(v bool) IS
	False(v bool) IS
	TrueM(v bool, msg string) IS
	FalseM(v bool, msg string) IS
	Fail(msg interface{})
}
