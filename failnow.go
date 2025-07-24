package gotestingmock

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/newmo-oss/go-caller"
)

// InvalidFailNowError is panicked when FailNow or related methods are called
// from a different goroutine than the test goroutine.
type InvalidFailNowError struct {
	File     string
	Line     int
	TestName string
	Method   string // FailNow, Fatal, Fatalf
	Format   string // the argument of Fatalf
	Args     []any  // the argument Fatal or Fatalf
}

func (err *InvalidFailNowError) Error() string {
	var method string
	switch err.Method {
	case "FailNow":
		method = "t.FailNow()"
	case "Fatal":
		method = fmt.Sprintf("t.Fatal(%q)", err.Args)
	case "Fatalf":
		method = fmt.Sprintf("t.Fatalf(%q, %v)", err.Format, err.Args)
	}
	return fmt.Sprintf("must not call %s on another goroutine with test %s in %s:%d", method, err.TestName, err.File, err.Line)
}

// StrictFailNow wraps testing.TB and panics if FailNow or related methods
// are called from a different goroutine than the test goroutine.
// This helps detect incorrect usage of t.FailNow in concurrent tests.
//
// Example:
//
//	func Test(t *testing.T) {
//		tb := gotestingmock.StrictFailNow(t)
//
//		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			if r.Method != http.MethodPost {
//				// This will panic because it's called from the handler's goroutine
//				tb.Fatal("unexpected method", r.Method)
//			}
//		}))
//		defer s.Close()
//
//		resp, err := http.Get(s.URL)
//		if err != nil {
//			tb.Fatal(err)  // This is safe - called from test goroutine
//		}
//		defer resp.Body.Close()
//	}
func StrictFailNow(tb testing.TB) *TB {
	testGID := goroutineID()

	var (
		file string
		line int
	)
	stacktrace := caller.New(1)
	if len(stacktrace) > 0 {
		file = stacktrace[0].File()
		line = stacktrace[0].Line()
	}

	return &TB{
		TB: tb,
		FailNowFunc: func() {
			if goroutineID() != testGID {
				panic(&InvalidFailNowError{
					File:     file,
					Line:     line,
					TestName: tb.Name(),
					Method:   "FailNow",
				})
			}
			tb.FailNow()
		},
		FatalFunc: func(args ...any) {
			if goroutineID() != testGID {
				panic(&InvalidFailNowError{
					File:     file,
					Line:     line,
					TestName: tb.Name(),
					Method:   "Fatal",
					Args:     args,
				})
			}
			tb.Fatal(args...)
		},
		FatalfFunc: func(format string, args ...any) {
			if goroutineID() != testGID {
				panic(&InvalidFailNowError{
					File:     file,
					Line:     line,
					TestName: tb.Name(),
					Method:   "Fatalf",
					Format:   format,
					Args:     args,
				})
			}
			tb.Fatalf(format, args...)
		},
	}
}

// goroutineID extracts the goroutine ID from the stack trace.
// It parses the runtime.Stack output which starts with "goroutine <id> ...".
func goroutineID() string {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	// 10: len("goroutine ")
	for i := 10; i < n; i++ {
		if buf[i] == ' ' {
			return string(buf[10:i])
		}
	}
	panic("cannot parse goroutine id from runtime.Stack")
}
