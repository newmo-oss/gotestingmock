package gotestingmock

import (
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

	testing.TB // for private method and unsupport method
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
	if tb.CleanupFunc != nil {
		tb.CleanupFunc(f)
	}
}

func (tb *TB) Error(args ...any) {
	tb.record.Failed = true
	if tb.ErrorFunc != nil {
		tb.ErrorFunc(args...)
	}
}

func (tb *TB) Errorf(format string, args ...any) {
	tb.record.Failed = true
	if tb.ErrorfFunc != nil {
		tb.ErrorfFunc(format, args...)
	}
}

func (tb *TB) Fail() {
	tb.record.Failed = true
	if tb.FailFunc != nil {
		tb.FailFunc()
	}
}

func (tb *TB) FailNow() {
	tb.record.Failed = true
	tb.record.Goexit = true
	if tb.FailNowFunc != nil {
		tb.FailNowFunc()
	} else {
		runtime.Goexit()
	}
}

func (tb *TB) Failed() bool {
	if tb.FailedFunc != nil {
		return tb.FailedFunc()
	}
	return tb.record.Failed
}

func (tb *TB) Fatal(args ...any) {
	tb.record.Failed = true
	tb.record.Goexit = true
	if tb.FatalFunc != nil {
		tb.FatalFunc(args...)
	}
	runtime.Goexit()
}

func (tb *TB) Fatalf(format string, args ...any) {
	tb.record.Failed = true
	tb.record.Goexit = true
	if tb.FatalfFunc != nil {
		tb.FatalfFunc(format, args...)
	}
	runtime.Goexit()
}

func (tb *TB) Helper() {
	if tb.HelperFunc != nil {
		tb.HelperFunc()
	}
}

func (tb *TB) Log(args ...any) {
	if tb.LogFunc != nil {
		tb.LogFunc(args...)
	}
}

func (tb *TB) Logf(format string, args ...any) {
	if tb.LogfFunc != nil {
		tb.LogfFunc(format, args...)
	}
}

func (tb *TB) Name() string {
	if tb.NameFunc != nil {
		return tb.NameFunc()
	}
	return ""
}

func (tb *TB) Setenv(key, value string) {
	if tb.SetenvFunc != nil {
		tb.SetenvFunc(key, value)
	}
}

func (tb *TB) Skip(args ...any) {
	tb.record.Skipped = true
	if tb.SkipFunc != nil {
		tb.SkipFunc(args...)
	}
}

func (tb *TB) SkipNow() {
	tb.record.Skipped = true
	tb.record.Goexit = true
	if tb.SkipNowFunc != nil {
		tb.SkipNowFunc()
	}
	runtime.Goexit()
}

func (tb *TB) Skipf(format string, args ...any) {
	tb.record.Skipped = true
	if tb.SkipfFunc != nil {
		tb.SkipfFunc(format, args...)
	}
}

func (tb *TB) Skipped() bool {
	if tb.SkippedFunc != nil {
		return tb.SkippedFunc()
	}
	return tb.record.Skipped
}

func (tb *TB) TempDir() string {
	if tb.TempDirFunc != nil {
		return tb.TempDirFunc()
	}
	return ""
}
