package main

import (
	"encoding/csv"
	"io"
	"regexp"
	"strconv"
	"strings"
)

var FIELD = regexp.MustCompile(`(?:[^,\\]|\\.)*`)

func parseIndexesList(list string) ([]int, error) {
	fields := FIELD.FindAllString(list, -1)

	nums := make([]int, len(fields))
	for i, field := range fields {
		num, err := strconv.Atoi(field)
		if err != nil {
			return nil, err
		}
		nums[i] = num
	}
	return nums, nil
}

func parseHeadersList(list string) ([]string, error) {
	fields := FIELD.FindAllString(list, -1)
	for i := 0; i < len(fields); i++ {
		fields[i] = strings.Replace(fields[i], `\,`, `,`, -1)
	}
	return fields, nil
}

type CSVScanner struct {
	text          string
	err           error
	parsedHeaders bool
	selector      Selector
	reader        *csv.Reader
}

func NewCSVScanner(s Selector, r io.Reader) *CSVScanner {
	return &CSVScanner{
		selector: s,
		reader:   csv.NewReader(r),
	}
}

func (c *CSVScanner) Err() error {
	if c.err == io.EOF {
		return nil
	}
	return c.err
}

func (c *CSVScanner) Bytes() []byte {
	return []byte(c.text)
}

func (c *CSVScanner) Text() string {
	return c.text
}

func (c *CSVScanner) Scan() bool {
	if c.err != nil {
		return false
	}

	recode, err := c.reader.Read()
	if err != nil {
		c.err = err
		c.text = ""
		return false
	}

	if !c.parsedHeaders && c.selector.RequireHeader() {
		err = c.selector.ParseHeader(recode)
		if err != nil {
			c.err = err
			c.text = ""
			return false
		}
		c.parsedHeaders = true
		return c.Scan()
	}

	values, err := c.selector.Select(recode)
	if err != nil {
		c.err = err
		c.text = ""
		return false
	}
	c.text = strings.Join(values, "\t")

	return true
}
