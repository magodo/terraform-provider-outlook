package utils

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestToPtr(t *testing.T) {
	var slice []int
	cases := []struct {
		in  interface{}
		out interface{}
	}{
		{
			in:  1,
			out: Int(1),
		},
		{
			in:  "a",
			out: String("a"),
		},
		{
			in:  true,
			out: Bool(true),
		},
		{
			in:  slice,
			out: &slice,
		},
	}
	for idx, c := range cases {
		out := ToPtr(c.in)
		if !reflect.DeepEqual(out, c.out) {
			t.Fatalf("%d failed\nexpected:\n%s\nactual:\n%s\n", idx, spew.Sdump(c.out), spew.Sdump(out))
		}
	}

}

func TestSafeDeref(t *testing.T) {
	type T1 int
	cases := []struct {
		in  interface{}
		out interface{}
	}{
		{
			in:  Int(1),
			out: 1,
		},
		{
			in:  (*int)(nil),
			out: 0,
		},
		{
			in:  (*T1)(nil),
			out: T1(0),
		},
	}
	for idx, c := range cases {
		out := SafeDeref(c.in)
		if !reflect.DeepEqual(out, c.out) {
			t.Fatalf("%d failed\nexpected:\n%s\nactual:\n%s\n", idx, spew.Sdump(c.out), spew.Sdump(out))
		}
	}
}
