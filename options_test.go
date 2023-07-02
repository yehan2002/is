package is

import (
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
		opts := &options{}
		opts.apply(o...)
		newIs(it, opts).Equal(v1, v2, "Test")
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
