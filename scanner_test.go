package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

type ScanResult struct {
	scan  bool
	text  string
	isErr bool
}

func ppResults(results []ScanResult) string {
	lines := make([]string, len(results))
	for i := 0; i < len(results); i++ {
		lines[i] = fmt.Sprintf("%#v", results[i])
		lines[i] = strings.TrimPrefix(lines[i], "main.ScanResult")
	}
	return strings.Join(lines, "\n")
}

type ScanTest struct {
	selector Selector
	src      string
	result   []ScanResult
}

type DummyAll struct {
}

func (d *DummyAll) DropHeaders() bool {
	return false
}

func (d *DummyAll) ParseHeaders(headers []string) error {
	return nil
}

func (d *DummyAll) Select(recode []string) ([]string, error) {
	return recode, nil
}

var scanTests = []ScanTest{
	{
		selector: &DummyAll{},
		src: `
1,1,1
2,4,8
`[1:],
		result: []ScanResult{
			{scan: true, text: "1\t1\t1", isErr: false},
			{scan: true, text: "2\t4\t8", isErr: false},
			{scan: false, text: "", isErr: false},
		},
	},
	{
		selector: &DummyAll{},
		src: `
1,2,3
1,2,3,4,5,6
`[1:],
		result: []ScanResult{
			{scan: true, text: "1\t2\t3", isErr: false},
			{scan: false, text: "", isErr: true},
			{scan: false, text: "", isErr: true},
		},
	},
}

func TestScan(t *testing.T) {
	for _, test := range scanTests {
		r := strings.NewReader(test.src)
		c := NewCSVScanner(test.selector, r)

		expect := test.result
		actual := make([]ScanResult, len(test.result))
		for i := 0; i < len(test.result); i++ {
			scan := c.Scan()
			actual[i] = ScanResult{
				scan:  scan,
				text:  c.Text(),
				isErr: c.Err() != nil,
			}
		}
		if !reflect.DeepEqual(actual, expect) {
			t.Errorf("src:\n%sgot:\n%s\nwant:\n%s",
				test.src,
				ppResults(actual), ppResults(expect))
		}
	}
}
