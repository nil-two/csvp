package main

import (
	"encoding/csv"
	"io"
	"strings"
)

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

	if !c.parsedHeaders && c.selector.RequireHeaders() {
		err = c.selector.ParseHeaders(recode)
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
