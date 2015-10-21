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
		indexes: []int{9},
	},
	{
		list:    "120",
		indexes: []int{119},
	},

	// Multiple indexes
	{
		list:    "10,120",
		indexes: []int{9, 119},
	},
	{
		list:    "10,120,50",
		indexes: []int{9, 119, 49},
	},
	{
		list:    "3,2,1,0",
		indexes: []int{2, 1, 0, -1},
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
			t.Errorf("parseIndexesList(%q) = %v, want %v",
				test.list, actual, expect)
		}
	}
}
