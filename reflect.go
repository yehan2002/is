package is

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
)

//maxStackLength the maximum amount of stack frames to read
const maxStackLength = 50

// getComment gets the Go comment from the specified line
// in the specified file.
// https://github.com/matryer/is/blob/master/is.go
func getComment() (comment string, ok bool) {
	var path string
	var line int
	for i := 0; ; i++ {
		_, path, line, ok = runtime.Caller(i)
		if !ok {
			return
		}
		if strings.Contains(path, "github.com/yehan2002/is/") {
			continue
		}
		break
	}

	f, err := os.Open(path)
	if err != nil {
		return "", false
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	i := 1
	for s.Scan() {
		if i == line {
			text := s.Text()
			commentI := strings.Index(text, "// ")
			if commentI == -1 {
				return "", false // no comment
			}
			text = text[commentI+2:]
			text = strings.TrimSpace(text)
			return text, true
		}
		i++
	}
	return "", false
}

func getStack(skip int) string {
	stackBuf := make([]uintptr, maxStackLength)
	length := runtime.Callers(skip+1, stackBuf[:])
	stack := stackBuf[:length]

	trace := ""
	frames := runtime.CallersFrames(stack)
	for {
		frame, more := frames.Next()
		if !isInternal(frame.File) {
			trace = trace + fmt.Sprintf("\n%s() \n\t%s:%d", frame.Function, frame.File, frame.Line)
		}
		if !more {
			break
		}
	}
	return trace[1:]
}

func isInternal(v string) bool {
	c := strings.Contains
	return c(v, "/go/src/testing/testing.go") || c(v, "/go/src/reflect/value.go") || c(v, "github.com/yehan2002/is/is.go") || c(v, "go/src/runtime/")
}

func captureStdout() func() string {
	var r, w *os.File
	var err error
	if r, w, err = os.Pipe(); err == nil {
		stdoutMux.Lock()
		old := os.Stdout
		os.Stdout = w

		buf := &strings.Builder{}

		fin := make(chan interface{}, 1)

		go func() {
			if _, err = io.Copy(buf, r); err != nil && err != io.EOF && err != io.ErrClosedPipe {
			}
			fin <- nil
		}()
		return func() string {
			os.Stdout = old
			stdoutMux.Unlock()
			r.Close()
			w.Close()
			<-fin
			return buf.String()
		}
	}
	panic(err)
}
