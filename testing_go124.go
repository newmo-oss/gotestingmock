//go:build go1.24

package gotestingmock

import (
	"context"
)

func (tb *TB) Context() context.Context {
	switch {
	case tb.ContextFunc != nil:
		return tb.ContextFunc()
	case tb.TB != nil:
		return tb.TB.Context()
	default:
		return context.Background()
	}
}
