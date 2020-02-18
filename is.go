package is

import (
	"fmt"
	"reflect"
	"testing"
)

type failable func(t *testing.T, msg interface{}, test interface{}, comment bool)

// IS is a package for writing tests

//New creates a new test helper
func New(t *testing.T) IS {
	return &baseTest{t: t, fail: basicFailable}
}

//NoColor disables color
func NoColor() {
	noColorFlag = true
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
	MustPanicCall(funct interface{}, args ...interface{})

	Err(v1 interface{}) IS
	True(v bool) IS
	False(v bool) IS
	TrueM(v bool, msg string) IS
	FalseM(v bool, msg string) IS
	Fail(msg interface{})
}

type baseTest struct {
	t    *testing.T
	fail failable
}

//Equal tests if the given values are equal
func (i *baseTest) Equal(v1, v2 interface{}) IS {
	if eq, err := deepEqual(v1, v2); !eq {
		i.fail(i.t, nil, fmt.Sprint(err), true)
	}
	return i
}

//EqualM like `Equal` but with a message
func (i *baseTest) EqualM(v1, v2 interface{}, msg string) IS {
	if eq, err := deepEqual(v1, v2); !eq {
		i.fail(i.t, msg, fmt.Sprint(err), false)
	}
	return i
}

//NotEqual tests if the given values are not equal
func (i *baseTest) NotEqual(v1, v2 interface{}) IS {
	if reflect.DeepEqual(v1, v2) {
		i.fail(i.t, nil, fmt.Sprintf("%#v is equal to %#v", v1, v2), true)
	}
	return i
}

//NotEqualM like `NotEqual` but with a message
func (i *baseTest) NotEqualM(v1, v2 interface{}, msg string) IS {
	if reflect.DeepEqual(v1, v2) {
		i.fail(i.t, msg, fmt.Sprintf("%#v is equal to %#v", v1, v2), false)
	}
	return i
}

//NotNil tests if the given value is not nil
func (i *baseTest) NotNil(v1 interface{}) IS {
	if v1 == nil {
		i.fail(i.t, nil, "unexpected nil value", true)
	}
	return i
}

//NotNil tests if the given value is not nil
func (i *baseTest) NotNilM(v1 interface{}, msg string) IS {
	if v1 == nil {
		i.fail(i.t, msg, "unexpected nil value", false)
	}
	return i
}

//Nil tests if the given value is nil
func (i *baseTest) Nil(v1 interface{}) IS {
	if v1 != nil {
		i.fail(i.t, nil, "value is not nil", true)
	}
	return i
}

//Nil tests if the given value is nil
func (i *baseTest) NilM(v1 interface{}, msg string) IS {
	if v1 != nil {
		i.fail(i.t, msg, "value is not nil", false)
	}
	return i
}

//Err tests for errors
func (i *baseTest) Err(v1 interface{}) IS {
	if v1 != nil {
		i.fail(i.t, v1, "unexpected error", false)
	}
	return i
}

//True tests if the given expression is true
func (i *baseTest) True(v bool) IS {
	if !v {
		i.fail(i.t, nil, "expected value to be true", true)
	}
	return i
}

//True tests if the given expression is true
func (i *baseTest) TrueM(v bool, msg string) IS {
	if !v {
		i.fail(i.t, msg, "expected value to be true", false)
	}
	return i
}

//False tests if the given expression is false
func (i *baseTest) False(v bool) IS {
	if v {
		i.fail(i.t, nil, "expected value to be false", true)
	}
	return i
}

//False tests if the given expression is false
func (i *baseTest) FalseM(v bool, msg string) IS {
	if v {
		i.fail(i.t, msg, "expected value to be false", false)
	}
	return i
}

//MustPanic tests if the code panics
func (i *baseTest) MustPanic() {
	if r := recover(); r == nil {
		i.fail(i.t, nil, "expected a panic", true)
	}
}

func (i *baseTest) MustPanicCall(funct interface{}, args ...interface{}) {
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
func (i *baseTest) Fail(msg interface{}) {
	i.fail(i.t, msg, nil, false)
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
