package is

import (
	"errors"
	"fmt"
	"math"
	"testing"

	"github.com/yehan2002/is/v2/internal"
)

type exportTest struct {
	w int
}

type exportTest2 struct {
	V int
	e exportTest
}

func testEq(t *testing.T, v1, v2 interface{}, o ...Option) *internal.Test {
	t.Helper()
	result := internal.Run(func(it internal.T) {
		newIs(it, newOptions(o)).Equal(v1, v2, "Test")
	})
	return result
}

func TestOptionCmpAllUnexported(t *testing.T) {
	result := testEq(t, exportTest{w: 1}, exportTest{w: 13}, CmpAllUnexported())
	if !result.Failed {
		t.Fatal("Test did not fail when v1 != v2")
	}

	// test comparing nested values
	result = testEq(t,
		exportTest2{e: exportTest{w: 13}},
		exportTest2{e: exportTest{w: 12}},
		CmpAllUnexported())
	if !result.Failed {
		t.Fatal("Test did not fail when v1 != v2")
	}
}

func TestOptionCmpUnexported(t *testing.T) {

	result := testEq(t, exportTest{w: 1}, exportTest{w: 13}, CmpUnexported(exportTest{}))
	if !result.Failed {
		t.Fatal("Test did not fail when v1 != v2")
	}

	result = testEq(t,
		exportTest2{e: exportTest{w: 13}},
		exportTest2{e: exportTest{w: 12}},
		CmpUnexported(exportTest{}))
	if result.Failed {
		t.Error(result.FailMessage)
		t.Fatal("Test failed because ignored fields were not equal")
	}

	result = testEq(t,
		exportTest2{e: exportTest{w: 13}},
		exportTest2{e: exportTest{w: 12}},
		CmpUnexported(exportTest2{}))
	if result.Failed {
		t.Error(result.FailMessage)
		t.Fatal("Test failed because ignored fields were not equal")
	}

	result = testEq(t,
		exportTest2{e: exportTest{w: 13}},
		exportTest2{e: exportTest{w: 12}},
		CmpUnexported(exportTest2{}, exportTest{}))
	if !result.Failed {
		t.Error(result.FailMessage)
		t.Fatal("Test did not fail when v1 != v2")
	}

}

func TestOptionEquateNaN(t *testing.T) {
	result := testEq(t, math.NaN(), math.NaN(), EquateNaN(false))
	if !result.Failed {
		t.Fatal("NaN values were considered to be equal without EquateNaNs")
	}

	result = testEq(t, math.NaN(), math.NaN(), EquateNaN(true))
	if result.Failed {
		t.Fatal("NaN values were not considered to be equal with EquateNaNs")
	}
}

func TestOptionEquateEmpty(t *testing.T) {
	result := testEq(t, []byte(nil), []byte{}, EquateEmpty(false))
	if !result.Failed {
		t.Fatal("[]byte(nil) and []byte{} should not be equal")
	}

	result = testEq(t, []byte(nil), []byte{}, EquateEmpty(true))
	if result.Failed {
		t.Fatal("[]byte(nil) and []byte{} were not considered to be equal with EquateEmpty")
	}
}

func TestOptionEquateErrors(t *testing.T) {
	err1 := errors.New("err1")
	err2 := fmt.Errorf("err2: %w", err1)

	result := testEq(t, err1, err2, EquateErrors(false))
	if !result.Failed {
		t.Fatal("err1 and err2 should not be equal")
	}

	result = testEq(t, err1, err2, EquateErrors(true))
	if result.Failed {
		t.Fatal("err1 and err2 were not considered to be equal with EquateErrors")
	}
}
