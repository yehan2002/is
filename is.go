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
	if !reflect.DeepEqual(v1, v2) {
		i.fail(i.t, nil, fmt.Sprintf("%#v is not equal to %#v", v1, v2), true)
	}
	return i
}

//EqualM like `Equal` but with a message
func (i *baseTest) EqualM(v1, v2 interface{}, msg string) IS {
	if !reflect.DeepEqual(v1, v2) {
		i.fail(i.t, msg, fmt.Sprintf("%#v is not equal to %#v", v1, v2), false)
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

//Fail fail the test immediately
func (i *baseTest) Fail(msg interface{}) {
	i.fail(i.t, msg, "", false)
}

func basicFailable(t *testing.T, msg interface{}, test interface{}, comment bool) {
	if msg != nil {
		fmt.Printf("--- FAIL: %s\n", msg)
	}
	if comment {
		if c, ok := getComment(); ok {
			fmt.Printf("--- FAIL: %s\n", c)
		}
	}
	if test != nil {
		fmt.Printf("--- Error: %s\n", test)
	}
	fmt.Println(getStack(3))
	t.FailNow()
}
