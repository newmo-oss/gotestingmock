package gotestingmock_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/newmo-oss/gotestingmock"
)

func TestStrictFailNow(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		methodCall func(tb testing.TB)
	}{
		"FailNow": {methodCall: func(tb testing.TB) { tb.FailNow() }},
		"Fatal":   {methodCall: func(tb testing.TB) { tb.Fatal("error") }},
		"Fatalf":  {methodCall: func(tb testing.TB) { tb.Fatalf("error") }},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			tb := gotestingmock.StrictFailNow(t)

			var recovered any
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				defer func() {
					recovered = recover()
				}()

				// This will panic because it's called from the handler's goroutine
				tt.methodCall(tb)
			}))
			defer s.Close()

			resp, err := http.Get(s.URL)
			if err != nil {
				tb.Fatal(err) // This is safe - called from test goroutine
			}
			defer resp.Body.Close()

			if recovered == nil {
				t.Error("expected panic did not occur")
			} else {
				err, ok := recovered.(*gotestingmock.InvalidFailNowError)
				if !ok {
					t.Error("unexpected panic:", recovered)
				} else {
					t.Log("expected panic:", err)
				}
			}
		})
	}

	t.Run("no panic", func(t *testing.T) {
		t.Parallel()

		record := gotestingmock.Run(func(tb *gotestingmock.TB) {
			tb = gotestingmock.StrictFailNow(tb)
			tb.Fatal("error")
		})

		if record.PanicValue != nil {
			t.Error("unexpected panic:", record.PanicValue)
		}

		if !record.Failed || !record.Goexit {
			t.Error("expected t.Fatal behavior did not occur")
		}
	})
}
