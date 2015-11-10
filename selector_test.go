package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func toLines(cells [][]string) string {
	lines := make([]string, len(cells))
	for i, line := range cells {
		lines[i] = fmt.Sprintf("%q", line)
	}
	return strings.Join(lines, "\n")
}

var selectAllTests = []struct {
	src [][]string
	dst [][]string
}{
	{
		src: [][]string{},
		dst: [][]string{},
	},
	{
		src: [][]string{
			{"aaa", "bbb", "ccc"},
			{"ddd", "eee", "fff"},
		},
		dst: [][]string{
			{"aaa", "bbb", "ccc"},
			{"ddd", "eee", "fff"},
		},
	},
	{
		src: [][]string{
			{"a", "bb", "ccc", "dddd", "eeeee"},
			{"f", "gg", "hhh", "iiii", "jjjjj"},
			{"k", "ll", "mmm", "nnnn", "kkkkk"},
		},
		dst: [][]string{
			{"a", "bb", "ccc", "dddd", "eeeee"},
			{"f", "gg", "hhh", "iiii", "jjjjj"},
			{"k", "ll", "mmm", "nnnn", "kkkkk"},
		},
	},
}

func TestSelectAll(t *testing.T) {
	a := NewAll()
	for _, test := range selectAllTests {
		expect := test.dst
		actual := make([][]string, len(test.src))
		for i, line := range test.src {
			actual[i], _ = a.Select(line)
		}
		if !reflect.DeepEqual(actual, expect) {
			t.Errorf("All:\nsrc:\n%s\ngot:\n%s\nwant:\n%s",
				toLines(test.src),
				toLines(actual), toLines(expect))
		}
	}
}

var newIndexesTests = []struct {
	list    string
	headers []string
	wantErr bool
	indexes []int
}{
	{
		list:    "",
		headers: []string{"", "", ""},
		indexes: []int{},
	},
	{
		list:    "1",
		headers: []string{"", "", ""},
		indexes: []int{0},
	},
	{
		list:    "3,1,4",
		headers: []string{"", "", ""},
		indexes: []int{2, 0, 3},
	},
	{
		list:    "2-4",
		headers: []string{"", "", "", "", ""},
		indexes: []int{1, 2, 3},
	},
	{
		list:    "2-8",
		headers: []string{"", "", "", "", ""},
		indexes: []int{1, 2, 3, 4},
	},
	{
		list:    "-2",
		headers: []string{"", "", "", "", ""},
		indexes: []int{0, 1},
	},
	{
		list:    "-8",
		headers: []string{"", "", "", "", ""},
		indexes: []int{0, 1, 2, 3, 4},
	},
	{
		list:    "2-",
		headers: []string{"", "", "", "", ""},
		indexes: []int{1, 2, 3, 4},
	},
	{
		list:    "8-",
		headers: []string{"", "", "", "", ""},
		indexes: []int{},
	},
	{
		list:    "0,5",
		headers: []string{"", "", ""},
		wantErr: true,
	},
	{
		list:    ",,",
		headers: []string{"", "", ""},
		wantErr: true,
	},
	{
		list:    "--,5",
		headers: []string{"", "", ""},
		wantErr: true,
	},
	{
		list:    "foo,5",
		headers: []string{"", "", ""},
		wantErr: true,
	},
	{
		list:    "1\\,5",
		headers: []string{"", "", ""},
		wantErr: true,
	},
}

func TestNewIndexes(t *testing.T) {
	for _, test := range newIndexesTests {
		i, err := NewIndexes(test.list)
		if err != nil {
			t.Errorf("NewIndexes(%q) returns %q, want nil",
				test.list, err)
			continue
		}
		switch {
		case test.wantErr:
			if err = i.ParseHeaders(test.headers); err == nil {
				t.Errorf("NewIndexes(%q).ParseHeaders(%q) returns nil, want err",
					test.list, test.headers)
			}
		default:
			if err = i.ParseHeaders(test.headers); err != nil {
				t.Errorf("NewIndexes(%q).ParseHeaders(%q) returns %q, want nil",
					test.list, test.headers, err)
				continue
			}
			expect := test.indexes
			actual := i.indexes
			if !reflect.DeepEqual(actual, expect) {
				t.Errorf("NewIndexes(%q).ParseHeaders(%q) = %v, want %v",
					test.list, test.headers, actual, expect)
			}
		}
	}
}

var selectIndexesTests = []struct {
	description string
	list        string
	headers     []string
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
		headers:     []string{"", "", ""},
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
		headers:     []string{"", "", ""},
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
		headers:     []string{"", "", ""},
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
		headers:     []string{"", "", ""},
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
		headers:     []string{"", "", ""},
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
			continue
		}
		if err := i.ParseHeaders(test.headers); err != nil {
			t.Errorf("NewIndexes(%q).ParseHeaders(%q) returns %q, want nil",
				test.list, test.headers, err)
			continue
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
			t.Errorf("%s:\nsrc:\n%s\ngot:\n%s\nwant:\n%s",
				self, toLines(test.src),
				toLines(actual), toLines(expect))
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
			t.Errorf("NewHeaders(%q).headers:\ngot :%q\nwant:%q",
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
			t.Errorf("%q.ParseHeaders(%q).indexes:\ngot :%v\nwant:%v",
				test.list, test.headers, actual, expect)
		}
	}
}

var selectHeadersTests = []struct {
	list    string
	headers []string
	src     [][]string
	dst     [][]string
}{
	{
		list:    "",
		headers: []string{"name", "price", "quantity"},
		src: [][]string{
			{"Apple", "60", "20"},
			{"Grapes", "140", "8"},
			{"Pineapple", "400", "2"},
			{"Orange", "50", "14"},
		},
		dst: [][]string{
			{},
			{},
			{},
			{},
		},
	},
	{
		list:    "name",
		headers: []string{"name", "price", "quantity"},
		src: [][]string{
			{"Apple", "60", "20"},
			{"Grapes", "140", "8"},
			{"Pineapple", "400", "2"},
			{"Orange", "50", "14"},
		},
		dst: [][]string{
			{"Apple"},
			{"Grapes"},
			{"Pineapple"},
			{"Orange"},
		},
	},
	{
		list:    "price,name",
		headers: []string{"name", "price", "quantity"},
		src: [][]string{
			{"Apple", "60", "20"},
			{"Grapes", "140", "8"},
			{"Pineapple", "400", "2"},
			{"Orange", "50", "14"},
		},
		dst: [][]string{
			{"60", "Apple"},
			{"140", "Grapes"},
			{"400", "Pineapple"},
			{"50", "Orange"},
		},
	},
	{
		list:    "quantity,quantity",
		headers: []string{"name", "price", "quantity"},
		src: [][]string{
			{"Apple", "60", "20"},
			{"Grapes", "140", "8"},
			{"Pineapple", "400", "2"},
			{"Orange", "50", "14"},
		},
		dst: [][]string{
			{"20", "20"},
			{"8", "8"},
			{"2", "2"},
			{"14", "14"},
		},
	},
	{
		list:    "date,name",
		headers: []string{"name", "price", "quantity"},
		src: [][]string{
			{"Apple", "60", "20"},
			{"Grapes", "140", "8"},
			{"Pineapple", "400", "2"},
			{"Orange", "50", "14"},
		},
		dst: [][]string{
			{"", "Apple"},
			{"", "Grapes"},
			{"", "Pineapple"},
			{"", "Orange"},
		},
	},
	{
		list:    "date,name,name,quantity,per,per",
		headers: []string{"name", "price", "quantity"},
		src: [][]string{
			{"Apple", "60", "20"},
			{"Grapes", "140", "8"},
			{"Pineapple", "400", "2"},
			{"Orange", "50", "14"},
		},
		dst: [][]string{
			{"", "Apple", "Apple", "20", "", ""},
			{"", "Grapes", "Grapes", "8", "", ""},
			{"", "Pineapple", "Pineapple", "2", "", ""},
			{"", "Orange", "Orange", "14", "", ""},
		},
	},
}

func TestSelectHeaders(t *testing.T) {
	for _, test := range selectHeadersTests {
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
		self := fmt.Sprintf("{list=%q, headers=%q}",
			test.list, test.headers)

		expect := test.dst
		actual := make([][]string, len(test.src))
		for i, line := range test.src {
			actual[i], err = h.Select(line)
			if err != nil {
				t.Errorf("%s.Select(%q) returns %q, want nil",
					self, line, err)
			}
		}
		if !reflect.DeepEqual(actual, expect) {
			t.Errorf("%s:\nsrc:\n%s\ngot:\n%s\nwant:\n%s",
				self, toLines(test.src),
				toLines(actual), toLines(expect))
		}
	}
}
