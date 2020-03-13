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
	//Equal tests if the given values are equal.
	//Struct fields with the tag `is:"-"` are ignored
	Equal(v1, v2 interface{}) IS
	//EqualM like `Equal` but with a message
	EqualM(v1, v2 interface{}, msg string) IS
	//NotEqual tests if the given values are not equal
	NotEqual(v1, v2 interface{}) IS
	//NotEqualM like `NotEqual` but with a message
	NotEqualM(v1, v2 interface{}, msg string) IS
	//NotNil tests if the given value is not nil
	NotNil(v1 interface{}) IS
	//NotNil tests if the given value is not nil
	NotNilM(v1 interface{}, msg string) IS
	//Nil tests if the given value is nil
	Nil(v1 interface{}) IS
	//Nil like `Nil` but with a message
	NilM(v1 interface{}, msg string) IS

	//MustPanic tests if the code panics
	MustPanic()
	//MustCallPanic tests if calling p will panic
	MustPanicCall(panicable)
	//MustPanicCallReflect tests if the calling the function will panic
	MustPanicCallReflect(funct interface{}, args ...interface{})

	//Err tests for errors
	Err(v1 interface{}) IS
	//True tests if the given expression is true
	True(v bool) IS
	//False tests if the given expression is false
	False(v bool) IS
	//True like `True` but with a message
	TrueM(v bool, msg string) IS
	//FalseM like `False` but with a message
	FalseM(v bool, msg string) IS
	//Fail fail the test immediately
	Fail(msg interface{})
}
