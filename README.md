# Is

A lightweight testing framework for golang.

## Usage

### Basic usage

```golang
func TestLoader(t *testing.T){
    l := loader{url: "http://example.com"}

    is := is.New(t)
    is.Equal(l.url,"http://example.com") // test default url

    r,err := l.Get()
    is.Err(err)

    is.Equal(r, testData) // the page content must match

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
    l.loader = loader{}
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

## API

* IS.Equal - Fails if the provided values are not are deeply equal
* IS.NotEqual - Fails if the provided values are  deeply equal
* IS.NotNil - Fails if the provided value is nil
* IS.Nil - Fails if the provided value is not nil
* IS.Err - Fails if the error is not nil
* IS.True - Fails if the provided value is not `true`
* IS.False - Fails if the provided value is not `false`
* IS.MustPanic - Fails if `recover()` returns nil
* IS.MustPanicCall - Calls the given function and fails if it does not panic
* IS.MustPanicCallReflect - Calls the given function with the given args and fails if it does not panic
* IS.Fail - Fails the test with the given message


## Color

`Is` defaults to not using color unless the `COLOR_TEST` is set to `true`
