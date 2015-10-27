package main

import (
	"fmt"
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

var newIndexesTests = []struct {
	description string
	list        string
	wantErr     bool
	indexes     []int
}{
	{
		list:    "",
		indexes: []int{},
	},
	{
		list:    "1",
		indexes: []int{0},
	},
	{
		list:    "3,1,4",
		indexes: []int{2, 0, 3},
	},
	{
		list:    "0,5",
		wantErr: true,
	},
	{
		list:    "-8,5",
		wantErr: true,
	},
	{
		list:    "foo,5",
		wantErr: true,
	},
	{
		list:    "1\\,5",
		wantErr: true,
	},
}

func TestNewIndexes(t *testing.T) {
	for _, test := range newIndexesTests {
		switch {
		case test.wantErr:
			_, err := NewIndexes(test.list)
			if err == nil {
				t.Errorf("NewIndexes(%q) returns nil, want err",
					test.list)
			}
		default:
			i, err := NewIndexes(test.list)
			if err != nil {
				t.Errorf("NewIndexes(%q) returns %q, want nil",
					test.list, err)
				continue
			}
			expect := test.indexes
			actual := i.indexes
			if !reflect.DeepEqual(actual, expect) {
				t.Errorf("NewIndexes(%q) = %v, want %v",
					test.list, actual, expect)
			}
		}
	}
}

var selectIndexesTests = []struct {
	description string
	list        string
	src         [][]string
	dst         [][]string
}{
	{
		description: "no input",
		list:        "1",
		src:         [][]string{},
		dst:         [][]string{},
	},
	{
		description: "only one index",
		list:        "1",
		src: [][]string{
			{"aaa", "bbb", "ccc"},
			{"ddd", "eee", "fff"},
		},
		dst: [][]string{
			{"aaa"},
			{"ddd"},
		},
	},
	{
		description: "index out of bounds",
		list:        "4",
		src: [][]string{
			{"aaa", "bbb", "ccc"},
			{"ddd", "eee", "fff"},
		},
		dst: [][]string{
			{""},
			{""},
		},
	},
	{
		description: "multiple indexes",
		list:        "3,1",
		src: [][]string{
			{"aaa", "bbb", "ccc"},
			{"ddd", "eee", "fff"},
		},
		dst: [][]string{
			{"ccc", "aaa"},
			{"fff", "ddd"},
		},
	},
	{
		description: "duplicated indexes",
		list:        "2,2,2",
		src: [][]string{
			{"aaa", "bbb", "ccc"},
			{"ddd", "eee", "fff"},
		},
		dst: [][]string{
			{"bbb", "bbb", "bbb"},
			{"eee", "eee", "eee"},
		},
	},
	{
		description: "battery",
		list:        "8,8,2,1,1,4",
		src: [][]string{
			{"a", "bb", "ccc", "dddd", "eeeee"},
			{"f", "gg", "hhh", "iiii", "jjjjj"},
			{"j", "kk", "lll", "mmmm", "nnnnn"},
		},
		dst: [][]string{
			{"", "", "bb", "a", "a", "dddd"},
			{"", "", "gg", "f", "f", "iiii"},
			{"", "", "kk", "j", "j", "mmmm"},
		},
	},
}

func TestSelectIndexes(t *testing.T) {
	for _, test := range selectIndexesTests {
		i, err := NewIndexes(test.list)
		if err != nil {
			t.Errorf("NewIndexes(%q) returns %q, want nil",
				test.list, err)
		}
		self := fmt.Sprintf("{list=%q, description=%q}",
			test.list, test.description)

		expect := test.dst
		actual := make([][]string, len(test.src))
		for j, line := range test.src {
			actual[j], err = i.Select(line)
			if err != nil {
				t.Errorf("%s: Select(%q) returns %q, want nil",
					self, line, err)
			}
		}
		if !reflect.DeepEqual(actual, expect) {
			t.Errorf("%s: %q: got %q, want %q",
				self, test.src, actual, expect)
		}
	}
}

var newHeadersTests = []struct {
	list    string
	headers []string
}{
	{
		list:    "",
		headers: []string{},
	},
	{
		list:    "name",
		headers: []string{"name"},
	},
	{
		list:    "name,price,quantity",
		headers: []string{"name", "price", "quantity"},
	},
	{
		list:    "a\\,b\\,c,d\\,e\\,f",
		headers: []string{"a,b,c", "d,e,f"},
	},
}

func TestNewHeaders(t *testing.T) {
	for _, test := range newHeadersTests {
		h, err := NewHeaders(test.list)
		if err != nil {
			t.Errorf("NewHeaders(%q) returns %q, want nil",
				test.list, err)
			continue
		}
		expect := test.headers
		actual := h.headers
		if !reflect.DeepEqual(actual, expect) {
			t.Errorf("NewHeaders(%q) = %v, want %v",
				test.list, actual, expect)
		}
	}
}

var headersParseHeadersTests = []struct {
	list    string
	headers []string
	indexes []int
}{
	{
		list:    "",
		headers: []string{"name", "price", "quantity"},
		indexes: []int{},
	},
	{
		list:    "name",
		headers: []string{"name", "price", "quantity"},
		indexes: []int{0},
	},
	{
		list:    "price,name",
		headers: []string{"name", "price", "quantity"},
		indexes: []int{1, 0},
	},
	{
		list:    "quantity,quantity",
		headers: []string{"name", "price", "quantity"},
		indexes: []int{2, 2},
	},
	{
		list:    "date,name",
		headers: []string{"name", "price", "quantity"},
		indexes: []int{-1, 0},
	},
	{
		list:    "date,name,name,quantity,per,per",
		headers: []string{"name", "price", "quantity"},
		indexes: []int{-1, 0, 0, 2, -1, -1},
	},
}

func TestHeadersParseHeaders(t *testing.T) {
	for _, test := range headersParseHeadersTests {
		h, err := NewHeaders(test.list)
		if err != nil {
			t.Errorf("NewHeaders(%q) returns %q, want nil",
				test.list, err)
			continue
		}
		if err = h.ParseHeaders(test.headers); err != nil {
			t.Errorf("%q.ParseHeaders(%q) returns %q, want nil",
				test.list, test.headers, err)
			continue
		}

		expect := test.indexes
		actual := h.indexes
		if !reflect.DeepEqual(actual, expect) {
			t.Errorf("%q.ParseHeaders(%q).indexes = %v, want %v",
				test.list, test.indexes, actual, expect)
		}
	}
}
