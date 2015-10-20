package main

import (
	"reflect"
	"testing"
)

var parseIndexesListTests = []struct {
	list    string
	indexes []int
}{
	// Only one index
	{
		list:    "10",
		indexes: []int{10},
	},
	{
		list:    "120",
		indexes: []int{120},
	},
}

func TestParseIndexesList(t *testing.T) {
	for _, test := range parseIndexesListTests {
		expect := test.indexes
		actual, err := parseIndexesList(test.list)
		if err != nil {
			t.Errorf("parseIndexesList(%q) returns %q, want nil",
				test.list, err)
		}
		if !reflect.DeepEqual(actual, expect) {
			t.Error("parseIndexesList(%q) = %v, want %v",
				test.list, actual, expect)
		}
	}
}
