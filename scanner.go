package main

import (
	"encoding/csv"
	"io"
	"strings"
)

type CSVScanner struct {
	outputDelimiter string
	text            string
	err             error
	parsedHeaders   bool
	selector        Selector
	reader          *csv.Reader
}

func NewCSVScanner(s Selector, r io.Reader) *CSVScanner {
	return &CSVScanner{
		outputDelimiter: "\t",
		selector:        s,
		reader:          csv.NewReader(r),
	}
}

func (c *CSVScanner) SetOutputDelimiter(s string) {
	c.outputDelimiter = s
}

func (c *CSVScanner) SetDelimiter(r rune) {
	c.reader.Comma = r
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

	if !c.parsedHeaders {
		err = c.selector.ParseHeaders(recode)
		if err != nil {
			c.err = err
			c.text = ""
			return false
		}
		c.parsedHeaders = true

		if c.selector.DropHeaders() {
			return c.Scan()
		}
	}

	recode, err = c.selector.Select(recode)
	if err != nil {
		c.err = err
		c.text = ""
		return false
	}
	c.text = strings.Join(recode, c.outputDelimiter)

	return true
}
