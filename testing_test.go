package gotestingmock_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/newmo-oss/gotestingmock"
)

func TestTB_Fields(t *testing.T) {
	t.Parallel()

	typ := reflect.TypeFor[gotestingmock.TB]()
	for i := range typ.NumField() {
		ft := typ.Field(i)
		if !strings.HasSuffix(ft.Name, "Func") &&
			ft.Type.Kind() != reflect.Func {
			continue
		}

		t.Run(ft.Name, func(t *testing.T) {
			t.Parallel()

			methodName := strings.TrimSuffix(ft.Name, "Func")

			var (
				call bool
				skip bool
			)
			rec := gotestingmock.Run(func(tb *gotestingmock.TB) {
				v := reflect.ValueOf(tb)
				fv := v.Elem().Field(i)
				method := v.MethodByName(methodName)
				if !method.IsValid() {
					skip = true
					return
				}

				/*
					tb.XxxFunc = func() {
						call = true
					}
					tb.Xxx()
				*/

				fv.Set(reflect.MakeFunc(ft.Type, func([]reflect.Value) []reflect.Value {
					call = true
					ret := make([]reflect.Value, fv.Type().NumOut())
					for i := range fv.Type().NumOut() {
						ret[i] = reflect.New(fv.Type().Out(i)).Elem()
					}
					return ret
				}))

				callWithZeros(method)
			})

			if skip {
				t.SkipNow()
			}

			if rec.PanicValue != nil {
				t.Fatal("unexpected panic:", rec.PanicValue)
			}

			if !call {
				t.Errorf("Field %s did not call with %s", ft.Name, methodName)
			}
		})
	}
}

func TestTB_DefaultMethod(t *testing.T) {
	t.Parallel()

	typ := reflect.TypeFor[gotestingmock.TB]()
	for i := range typ.NumField() {
		ft := typ.Field(i)
		if !strings.HasSuffix(ft.Name, "Func") &&
			ft.Type.Kind() != reflect.Func {
			continue
		}

		t.Run(ft.Name, func(t *testing.T) {
			t.Parallel()

			methodName := strings.TrimSuffix(ft.Name, "Func")

			var (
				call bool
				skip bool
			)
			rec := gotestingmock.Run(func(parent *gotestingmock.TB) {
				tb := &gotestingmock.TB{TB: parent}

				pv := reflect.ValueOf(parent)
				pfv := pv.Elem().Field(i)

				v := reflect.ValueOf(tb)
				method := v.MethodByName(methodName)
				if !method.IsValid() {
					skip = true
					return
				}

				/*
					parent.XxxFunc = func() {
						call = true
					}
					parent.Xxx()
				*/

				pfv.Set(reflect.MakeFunc(ft.Type, func([]reflect.Value) []reflect.Value {
					call = true
					ret := make([]reflect.Value, pfv.Type().NumOut())
					for i := range pfv.Type().NumOut() {
						ret[i] = reflect.New(pfv.Type().Out(i)).Elem()
					}
					return ret
				}))

				callWithZeros(method)
			})

			if skip {
				t.SkipNow()
			}

			if rec.PanicValue != nil {
				t.Fatal("unexpected panic:", rec.PanicValue)
			}

			if !call {
				t.Errorf("(*gotestingmock.TB).%[1]s did not call with (testing.TB).%[1]s", methodName)
			}
		})
	}
}

func TestRecord(t *testing.T) {
	t.Parallel()

	cases := []struct {
		method      string
		wantSkipped bool
		wantFailed  bool
		wantGoexit  bool
	}{
		{"Error", false, true, false},
		{"Errorf", false, true, false},
		{"Fail", false, true, false},
		{"FailNow", false, true, true},
		{"Fatal", false, true, true},
		{"Fatalf", false, true, true},
		{"Skip", true, false, false},
		{"SkipNow", true, false, true},
		{"Skipf", true, false, false},
	}

	for _, tt := range cases {
		t.Run(tt.method, func(t *testing.T) {
			t.Parallel()
			rec := gotestingmock.Run(func(tb *gotestingmock.TB) {
				v := reflect.ValueOf(tb)
				callWithZeros(v.MethodByName(tt.method))
			})

			if rec.PanicValue != nil {
				t.Fatal("unexpected panic:", rec.PanicValue)
			}

			if got := rec.Skipped; got != tt.wantSkipped {
				t.Errorf("Skipped does not match: (got, want) = (%v, %v)", got, tt.wantSkipped)
			}

			if got := rec.Failed; got != tt.wantFailed {
				t.Errorf("Failed does not match: (got, want) = (%v, %v)", got, tt.wantFailed)
			}

			if got := rec.Goexit; got != tt.wantGoexit {
				t.Errorf("Goexit does not match: (got, want) = (%v, %v)", got, tt.wantGoexit)
			}
		})
	}
}

func TestRecord_PanicValue(t *testing.T) {
	t.Parallel()

	want := "panic"
	rec := gotestingmock.Run(func(tb *gotestingmock.TB) {
		panic(want)
	})

	if got := rec.PanicValue; got != want {
		t.Errorf("PanicValue does not match: (got, want) = (%v, %v)", got, want)
	}
}

func callWithZeros(v reflect.Value) {
	in := make([]reflect.Value, v.Type().NumIn())
	for i := range v.Type().NumIn() {
		in[i] = reflect.New(v.Type().In(i)).Elem()
	}
	v.Call(in)
}
