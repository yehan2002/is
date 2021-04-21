package is

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/yehan2002/is/internal"
)

type baseTest struct {
	t    *testing.T
	fail failable
}

//Equal tests if the given values are equal.
//Struct fields with the tag `is:"-"` are ignored
func (i *baseTest) Equal(v1, v2 interface{}, msg ...interface{}) IS {
	if eq, err := internal.IsEqual(v1, v2); !eq {
		i.fail(i.t, fmt.Sprint(err), true, msg)
	}
	return i
}

//NotEqual tests if the given values are not equal
func (i *baseTest) NotEqual(v1, v2 interface{}, msg ...interface{}) IS {
	if reflect.DeepEqual(v1, v2) {
		i.fail(i.t, fmt.Sprintf("%#v is equal to %#v", v1, v2), true, msg)
	}
	return i
}

//NotNil tests if the given value is not nil
func (i *baseTest) NotNil(v1 interface{}, msg ...interface{}) IS {
	if v1 == nil && reflect.ValueOf(v1).IsNil() {
		i.fail(i.t, "unexpected nil value", true, msg)
	}
	return i
}

//Nil tests if the given value is nil
func (i *baseTest) Nil(v1 interface{}, msg ...interface{}) IS {
	if v1 != nil && !reflect.ValueOf(v1).IsNil() {
		i.fail(i.t, "value is not nil", true, msg)
	}
	return i
}

//Err tests for errors
func (i *baseTest) Err(v1 interface{}, msg ...interface{}) IS {
	if v1 != nil {
		i.fail(i.t, "unexpected error", false, append([]interface{}{v1}, msg...))
	}
	return i
}

//True tests if the given expression is true
func (i *baseTest) True(v bool, msg ...interface{}) IS {
	if !v {
		i.fail(i.t, "expected value to be true", true, msg)
	}
	return i
}

//False tests if the given expression is false
func (i *baseTest) False(v bool, msg ...interface{}) IS {
	if v {
		i.fail(i.t, "expected value to be false", true, msg)
	}
	return i
}

//MustPanic tests if the code panics
func (i *baseTest) MustPanic(msg ...interface{}) {
	if r := recover(); r == nil {
		i.fail(i.t, "expected a panic", true, msg)
	}
}

//MustCallPanic tests if calling p will panic
func (i *baseTest) MustPanicCall(p panicable, msg ...interface{}) {
	defer i.MustPanic(msg)
	p()
}

//MustPanicCallReflect tests if the calling the function will panic
func (i *baseTest) MustPanicCallReflect(funct interface{}, args ...interface{}) {
	funcType := reflect.TypeOf(funct)
	if funcType.Kind() != reflect.Func {
		panic("`funct` is not a function")
	}
	if funcType.NumIn() != len(args) {
		panic("Invalid number of args")
	}

	rArgs := make([]reflect.Value, len(args))
	for i, arg := range args {
		rArgs[i] = reflect.ValueOf(arg)
		if !rArgs[i].IsValid() {
			panic("untyped nil provided as arg")
		}
		if funcType.In(i) != rArgs[i].Type() {
			panic("cannot assign " + funcType.In(i).String() + " to " + rArgs[i].Type().String())
		}
	}
	funcValue := reflect.ValueOf(funct)
	if !funcValue.IsValid() {
		panic("invalid function provided")
	}

	defer i.MustPanic()
	funcValue.Call(rArgs)
}

//Fail fail the test immediately
func (i *baseTest) Fail(msg ...interface{}) {
	i.fail(i.t, msg, false, msg)
}

//EqualM same as Equal.
func (i *baseTest) EqualM(v1, v2 interface{}, msg ...interface{}) IS {
	return i.Equal(v1, v2, msg...)
}

//NotEqual same as NotEqual.
func (i *baseTest) NotEqualM(v1, v2 interface{}, msg ...interface{}) IS {
	return i.NotEqual(v1, v2, msg...)
}

//NotNil same as NotNil.
func (i *baseTest) NotNilM(v1 interface{}, msg ...interface{}) IS { return i.NotNil(v1, msg...) }

//Nil same as Nil.
func (i *baseTest) NilM(v1 interface{}, msg ...interface{}) IS { return i.Nil(v1, msg...) }

//Err same as Err.
func (i *baseTest) ErrM(v1 interface{}, msg ...interface{}) IS { return i.Err(v1, msg...) }

//True same as True.
func (i *baseTest) TrueM(v bool, msg ...interface{}) IS { return i.True(v, msg...) }

//False same as False.
func (i *baseTest) FalseM(v bool, msg ...interface{}) IS { return i.False(v, msg...) }

func basicFailable(t *testing.T, test interface{}, comment bool, msg []interface{}) {
	if test != nil {
		messages.Err2.Print(true, test)
	}

	if msg != nil && len(msg) > 0 {
		messages.Err1.Print(true, fmt.Sprint(msg...))
	} else if comment {
		if c, ok := internal.GetComment(); ok {
			messages.Err1.Print(true, c)
		}
	}

	fmt.Println(internal.GetStack(3))
	t.FailNow()
}
