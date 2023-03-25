# Is [![Go](https://github.com/yehan2002/is/actions/workflows/go.yml/badge.svg)](https://github.com/yehan2002/is/actions/workflows/go.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/yehan2002/is/v2)](https://goreportcard.com/report/github.com/yehan2002/is/v2) [![codecov](https://codecov.io/gh/yehan2002/is/branch/v2/graph/badge.svg?token=I29DKX58GJ)](https://codecov.io/gh/yehan2002/is) [![Go Reference](https://pkg.go.dev/badge/github.com/yehan2002/is/v2.svg)](https://pkg.go.dev/github.com/yehan2002/is/v2)


A lightweight testing framework for golang.

## Usage

### Basic usage

```golang
func TestLoader(t *testing.T){
    is := is.New(t)
    
    l := loader{url: "http://example.com"}
    
    r, err := l.Get()
    is(l.url == "http://example.com", "calling Get() should not modify url")
    if err == nil{
        is(r != nil, "response should not be nil if err != nil")
        is.Equal(r, testData, "the page content must match")
    } else {
         is.Log("Failed to get test data. Skipping test.")
    }

}
}

```

### Test Suites

```golang
type LoaderTest struct{
    data []byte
    loader *loader
}

func (l *LoaderTest) Setup(){
    l.loader = &loader{}
}

func (l *LoaderTest) TestUrl(is is.Is){
    // tests go here
}

func (l *LoaderTest) Teardown(){
    l.loader.Close()
}

func TestLoader(t *testing.T){
    is.Suite(t, &LoaderTest{})
}

```

## Functions

* Is.Equal - Fails if the provided values are not are deeply equal
* Is.Panic - Fails if `recover()` returns nil
* Is.Fail - Fails the test with the given message
