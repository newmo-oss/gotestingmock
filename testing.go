package gotestingmock

import (
	"context"
	"io"
	"runtime"
	"sync"
	"testing"
)

// TB is mock for testing.TB.
// XxxFunc is a mock function for the method Xxx of testing.TB.
// If you confirm more usage please see [Run].
type TB struct {
	record Record

	// mock funcs
	CleanupFunc func(func())
	ErrorFunc   func(args ...any)
	ErrorfFunc  func(format string, args ...any)
	FailFunc    func()
	FailNowFunc func()
	FailedFunc  func() bool
	FatalFunc   func(args ...any)
	FatalfFunc  func(format string, args ...any)
	HelperFunc  func()
	LogFunc     func(args ...any)
	LogfFunc    func(format string, args ...any)
	NameFunc    func() string
	SetenvFunc  func(key, value string)
	SkipFunc    func(args ...any)
	SkipNowFunc func()
	SkipfFunc   func(format string, args ...any)
	SkippedFunc func() bool
	TempDirFunc func() string
	ContextFunc func() context.Context // for Go1.24
	OutputFunc  func() io.Writer       // for Go1.25

	testing.TB // for default behavior and private method
}

// Record records the result of [Run].
type Record struct {
	Failed     bool
	Skipped    bool
	Goexit     bool
	PanicValue any
}

// Run runs the given mocking test function with [testing.TB].
// The f can be described in the same way as the test function of the go test.
// Run call the f on new goroutine.
// The return value records whether the test function failed (e.g. t.Error),
// was skipped (e.g. t.Skip), failed and exited its goroutine (e.g. t.Fatal)
// or panic occured.
func Run(f func(*TB)) *Record {
	var (
		tb  TB
		wg  sync.WaitGroup
		ret *Record
	)
	wg.Add(1)
	go func() {
		defer func() {
			if p := recover(); p != nil {
				tb.record.PanicValue = p
			}
			record := tb.record // copy
			ret = &record
			wg.Done()
		}()
		f(&tb)
	}()
	wg.Wait()
	return ret
}

func (tb *TB) Cleanup(f func()) {
	switch {
	case tb.CleanupFunc != nil:
		tb.CleanupFunc(f)
	case tb.TB != nil:
		tb.TB.Cleanup(f)
	}
}

func (tb *TB) Error(args ...any) {
	tb.record.Failed = true
	switch {
	case tb.ErrorFunc != nil:
		tb.ErrorFunc(args...)
	case tb.TB != nil:
		tb.TB.Error(args...)
	}
}

func (tb *TB) Errorf(format string, args ...any) {
	tb.record.Failed = true
	switch {
	case tb.ErrorfFunc != nil:
		tb.ErrorfFunc(format, args...)
	case tb.TB != nil:
		tb.TB.Errorf(format, args...)
	}
}

func (tb *TB) Fail() {
	tb.record.Failed = true
	switch {
	case tb.FailFunc != nil:
		tb.FailFunc()
	case tb.TB != nil:
		tb.TB.Fail()
	}
}

func (tb *TB) FailNow() {
	tb.record.Failed = true
	tb.record.Goexit = true

	switch {
	case tb.FailNowFunc != nil:
		tb.FailNowFunc()
		runtime.Goexit()
	case tb.TB != nil:
		tb.TB.FailNow()
	default:
		runtime.Goexit()
	}
}

func (tb *TB) Failed() bool {
	switch {
	case tb.FailedFunc != nil:
		return tb.FailedFunc()
	case tb.TB != nil:
		return tb.TB.Failed()
	default:
		return tb.record.Failed
	}
}

func (tb *TB) Fatal(args ...any) {
	tb.record.Failed = true
	tb.record.Goexit = true
	switch {
	case tb.FatalFunc != nil:
		tb.FatalFunc(args...)
		runtime.Goexit()
	case tb.TB != nil:
		tb.TB.Fatal(args...)
	default:
		runtime.Goexit()
	}
}

func (tb *TB) Fatalf(format string, args ...any) {
	tb.record.Failed = true
	tb.record.Goexit = true

	switch {
	case tb.FatalfFunc != nil:
		tb.FatalfFunc(format, args...)
		runtime.Goexit()
	case tb.TB != nil:
		tb.TB.Fatalf(format, args...)
	default:
		runtime.Goexit()
	}
}

func (tb *TB) Helper() {
	switch {
	case tb.HelperFunc != nil:
		tb.HelperFunc()
	case tb.TB != nil:
		tb.TB.Helper()
	}
}

func (tb *TB) Log(args ...any) {
	switch {
	case tb.LogFunc != nil:
		tb.LogFunc(args...)
	case tb.TB != nil:
		tb.TB.Log(args...)
	}
}

func (tb *TB) Logf(format string, args ...any) {
	switch {
	case tb.LogfFunc != nil:
		tb.LogfFunc(format, args...)
	case tb.TB != nil:
		tb.TB.Logf(format, args...)
	}
}

func (tb *TB) Name() string {
	switch {
	case tb.NameFunc != nil:
		return tb.NameFunc()
	case tb.TB != nil:
		return tb.TB.Name()
	default:
		return ""
	}
}

func (tb *TB) Setenv(key, value string) {
	switch {
	case tb.SetenvFunc != nil:
		tb.SetenvFunc(key, value)
	case tb.TB != nil:
		tb.TB.Setenv(key, value)
	}
}

func (tb *TB) Skip(args ...any) {
	tb.record.Skipped = true
	switch {
	case tb.SkipFunc != nil:
		tb.SkipFunc(args...)
	case tb.TB != nil:
		tb.TB.Skip(args...)
	}
}

func (tb *TB) SkipNow() {
	tb.record.Skipped = true
	tb.record.Goexit = true
	switch {
	case tb.SkipNowFunc != nil:
		tb.SkipNowFunc()
		runtime.Goexit()
	case tb.TB != nil:
		tb.TB.SkipNow()
	default:
		runtime.Goexit()
	}
}

func (tb *TB) Skipf(format string, args ...any) {
	tb.record.Skipped = true
	switch {
	case tb.SkipfFunc != nil:
		tb.SkipfFunc(format, args...)
	case tb.TB != nil:
		tb.TB.Skipf(format, args...)
	}
}

func (tb *TB) Skipped() bool {
	switch {
	case tb.SkippedFunc != nil:
		return tb.SkippedFunc()
	case tb.TB != nil:
		return tb.TB.Skipped()
	default:
		return tb.record.Skipped
	}
}

func (tb *TB) TempDir() string {
	switch {
	case tb.TempDirFunc != nil:
		return tb.TempDirFunc()
	case tb.TB != nil:
		return tb.TB.TempDir()
	default:
		return ""
	}
}
