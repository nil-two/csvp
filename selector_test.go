package main

import (
	"reflect"
	"testing"
)

var selectAllTests = []struct {
	src []string
	dst []string
}{
	{
		src: []string{},
		dst: []string{},
	},
	{
		src: []string{"aaa", "bbb", "ccc"},
		dst: []string{"aaa", "bbb", "ccc"},
	},
	{
		src: []string{"a", "bb", "ccc", "dddd", "eeeee"},
		dst: []string{"a", "bb", "ccc", "dddd", "eeeee"},
	},
}

func TestSelectAll(t *testing.T) {
	a := NewAll()
	for _, test := range selectAllTests {
		expect := test.dst
		actual, err := a.Select(test.src)
		if err != nil {
			t.Errorf("All.Select(%q) returns %q, want nil",
				test.src, err)
		}
		if !reflect.DeepEqual(actual, expect) {
			t.Errorf("All.Select(%q) = %q, want %q",
				test.src, actual, expect)
		}
	}
}
