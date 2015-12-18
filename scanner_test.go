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
	dropHeaders              bool
	causeErrorAtParseHeaders bool
	causeErrorAtSelect       bool
}

func (d *DummyAll) DropHeaders() bool {
	return d.dropHeaders
}

func (d *DummyAll) ParseHeaders(headers []string) error {
	if d.causeErrorAtParseHeaders {
		return fmt.Errorf("parse error")
	}
	return nil
}

func (d *DummyAll) Select(recode []string) ([]string, error) {
	if d.causeErrorAtSelect {
		return nil, fmt.Errorf("select error")
	}
	return recode, nil
}

var scanTests = []ScanTest{
	{
		selector: &DummyAll{dropHeaders: false},
		src: `
`[1:],
		result: []ScanResult{
			{scan: false, text: "", isErr: false},
		},
	},
	{
		selector: &DummyAll{dropHeaders: false},
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
		selector: &DummyAll{dropHeaders: false},
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
	{
		selector: &DummyAll{dropHeaders: true},
		src: `
`[1:],
		result: []ScanResult{
			{scan: false, text: "", isErr: false},
		},
	},
	{
		selector: &DummyAll{dropHeaders: true},
		src: `
a,b,c
`[1:],
		result: []ScanResult{
			{scan: false, text: "", isErr: false},
		},
	},
	{
		selector: &DummyAll{dropHeaders: true},
		src: `
a,b,c
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
		selector: &DummyAll{dropHeaders: true},
		src: `
a,b,c
1,2,3
1,2,3,4,5,6
`[1:],
		result: []ScanResult{
			{scan: true, text: "1\t2\t3", isErr: false},
			{scan: false, text: "", isErr: true},
			{scan: false, text: "", isErr: true},
		},
	},
	{
		selector: &DummyAll{causeErrorAtParseHeaders: true},
		src: `
1,1,1
2,4,8
`[1:],
		result: []ScanResult{
			{scan: false, text: "", isErr: true},
			{scan: false, text: "", isErr: true},
			{scan: false, text: "", isErr: true},
		},
	},
	{
		selector: &DummyAll{causeErrorAtSelect: true},
		src: `
1,1,1
2,4,8
`[1:],
		result: []ScanResult{
			{scan: false, text: "", isErr: true},
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

func TestBytes(t *testing.T) {
	selector := &DummyAll{dropHeaders: false}
	r := strings.NewReader("100,200,300\n")
	c := NewCSVScanner(selector, r)
	c.Scan()

	expect := c.Text()
	actual := string(c.Bytes())
	if actual != expect {
		t.Errorf("expected string(c.Byte()) to equal c.Text()")
	}
}

var delimiterTests = []struct {
	delimiter rune
	src       string
	dst       []string
}{
	{
		delimiter: ' ',
		src: `
aaa bbb ccc
ddd eee fff
`[1:],
		dst: []string{
			"aaa\tbbb\tccc",
			"ddd\teee\tfff",
		},
	},
	{
		delimiter: '/',
		src: `
aaa/bbb/ccc/ddd
eee/fff/ggg/hhh
`[1:],
		dst: []string{
			"aaa\tbbb\tccc\tddd",
			"eee\tfff\tggg\thhh",
		},
	},
}

func TestScanWithDelimiter(t *testing.T) {
	selector := &DummyAll{}
	for _, test := range delimiterTests {
		r := strings.NewReader(test.src)
		c := NewCSVScanner(selector, r)
		c.SetDelimiter(test.delimiter)

		expect := test.dst
		actual := make([]string, 0)
		for c.Scan() {
			actual = append(actual, c.Text())
		}
		if c.Err() != nil {
			t.Errorf("delimiter:%q\nsrc:\n%sgot: %v, want nil",
				test.delimiter, test.src, c.Err())
			continue
		}
		if !reflect.DeepEqual(actual, expect) {
			t.Errorf("delimiter:%q\nsrc:\n%sgot:\n%s\nwant:\n%s",
				test.delimiter,
				test.src,
				strings.Join(actual, "\n"),
				strings.Join(expect, "\n"))
		}
	}
}

var outputDelimiterTests = []struct {
	outputDelimiter string
	src             string
	dst             []string
}{
	{
		outputDelimiter: "...",
		src: `
aaa,bbb,ccc
ddd,eee,fff
`[1:],
		dst: []string{
			"aaa...bbb...ccc",
			"ddd...eee...fff",
		},
	},
	{
		outputDelimiter: "→",
		src: `
aaa,bbb,ccc,ddd
eee,fff,ggg,hhh
`[1:],
		dst: []string{
			"aaa→bbb→ccc→ddd",
			"eee→fff→ggg→hhh",
		},
	},
}

func TestScanWithOutputDelimiter(t *testing.T) {
	selector := &DummyAll{}
	for _, test := range outputDelimiterTests {
		r := strings.NewReader(test.src)
		c := NewCSVScanner(selector, r)
		c.SetOutputDelimiter(test.outputDelimiter)

		expect := test.dst
		actual := make([]string, 0)
		for c.Scan() {
			actual = append(actual, c.Text())
		}
		if c.Err() != nil {
			t.Errorf("output-delimiter:%q\nsrc:\n%s"+
				"got: %v, want nil",
				test.outputDelimiter, test.src, c.Err())
			continue
		}
		if !reflect.DeepEqual(actual, expect) {
			t.Errorf("output-delimiter:%q\nsrc:\n%s"+
				"got:\n%s\nwant:\n%s",
				test.outputDelimiter,
				test.src,
				strings.Join(actual, "\n"),
				strings.Join(expect, "\n"))
		}
	}
}
