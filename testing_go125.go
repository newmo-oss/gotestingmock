//go:build go1.25

package gotestingmock

import (
	"io"
	"os"
)

func (tb *TB) Output() io.Writer {
	switch {
	case tb.OutputFunc != nil:
		return tb.OutputFunc()
	case tb.TB != nil:
		return tb.TB.Output()
	default:
		return os.Stdout
	}
}
