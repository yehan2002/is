//Package is provides helper functions for testing
package is

import (
	"reflect"
	"testing"

	"github.com/yehan2002/is/internal"
)

//IS a helper for writing tests.
type IS interface {
	//Equal tests if the given values are equal.
	//Struct fields with the tag `is:"-"` are ignored
	Equal(v1, v2 interface{}, msg ...interface{}) IS
	//NotEqual tests if the given values are not equal
	NotEqual(v1, v2 interface{}, msg ...interface{}) IS
	//NotNil tests if the given value is not nil
	NotNil(v1 interface{}, msg ...interface{}) IS
	//Nil tests if the given value is nil
	Nil(v1 interface{}, msg ...interface{}) IS
	//Err tests for errors
	Err(v1 interface{}, msg ...interface{}) IS
	//True tests if the given expression is true
	True(v bool, msg ...interface{}) IS
	//False tests if the given expression is false
	False(v bool, msg ...interface{}) IS
	//Fail fail the test immediately
	Fail(msg ...interface{})

	//MustPanic tests if the code panics
	MustPanic(msg ...interface{})
	//MustCallPanic tests if calling p will panic
	MustPanicCall(p panicable, msg ...interface{})
	//MustPanicCallReflect tests if the calling the function will panic
	MustPanicCallReflect(funct interface{}, args ...interface{})

	//EqualM same as Equal.
	EqualM(v1, v2 interface{}, msg ...interface{}) IS
	//NotEqual same as NotEqual.
	NotEqualM(v1, v2 interface{}, msg ...interface{}) IS
	//NotNilM same as NotNil.
	NotNilM(v1 interface{}, msg ...interface{}) IS
	//NilM same as Nil.
	NilM(v1 interface{}, msg ...interface{}) IS
	//ErrM same as Err.
	ErrM(v1 interface{}, msg ...interface{}) IS
	//TrueM same as True.
	TrueM(v bool, msg ...interface{}) IS
	//FalseM same as False.
	FalseM(v bool, msg ...interface{}) IS
}

//New creates a new test helper
func New(t *testing.T) IS { return &baseTest{t: t, fail: basicFailable} }

//NoColor disables color
func NoColor() {
	internal.NoColorFlag = true
}

type panicable func()
type failable func(t *testing.T, test interface{}, comment bool, msg []interface{})

var messages = internal.Messages
var isType = reflect.TypeOf((*IS)(nil)).Elem()
