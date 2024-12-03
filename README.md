# gotestingmock [![Go Reference](https://pkg.go.dev/badge/github.com/newmo-oss/gotestingmocko.svg)](https://pkg.go.dev/github.com/newmo-oss/gotestingmock)[![Go Report Card](https://goreportcard.com/badge/github.com/newmo-oss/gotestingmock)](https://goreportcard.com/report/github.com/newmo-oss/gotestingmock)

gotestingmock mocking utilities for unit test in Go.

## Usage

```go
func Test(t *testing.T) {
	t.Parallel()

	// gotestingmock.Run simulates a test function on a goroutine.
	got := gotestingmock.Run(func(tb *gotestingmock.TB) {
		// The test helper can use *gotestingmock.TB as testing.TB
		// which is implemented  by testing.T, testing.B and testing.F.
		MyTestHelper(tb, "arg1")
	})

	// Check if the test helper failed with t.Error, t.Fatal or similar methods.
	if !got.Failed {
		t.Error("expected failed did not occur")
	}

	// Check that the test helper has panicked.
	if got.PanicValue != nil {
		t.Error("unexpected panic:", got.PanicValue)
	}
}
```

## License
MIT
