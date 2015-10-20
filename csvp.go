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

type CSVScanner struct {
	indexes []int
	text    string
	err     error
	reader  *csv.Reader
}

func NewCSVScanner(indexes []int, r io.Reader) *CSVScanner {
	return &CSVScanner{
		indexes: indexes,
		reader:  csv.NewReader(r),
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

	var values []string
	for _, index := range c.indexes {
		values = append(values, recode[index])
	}
	c.text = strings.Join(values, "\t")

	return true
}
