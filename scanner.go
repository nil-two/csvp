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

func (c *CSVScanner) SetDelimiter(ch rune) {
	c.reader.Comma = ch
}

func (c *CSVScanner) SetOutputDelimiter(s string) {
	c.outputDelimiter = s
}

func (c *CSVScanner) InitializeReader(r io.Reader) {
	ch := c.reader.Comma
	c.reader = csv.NewReader(r)
	c.reader.Comma = ch
	c.parsedHeaders = false
	c.err = nil
	c.text = ""
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

	record, err := c.reader.Read()
	if err != nil {
		c.err = err
		c.text = ""
		return false
	}

	if !c.parsedHeaders {
		err = c.selector.ParseHeaders(record)
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

	record, err = c.selector.Select(record)
	if err != nil {
		c.err = err
		c.text = ""
		return false
	}
	c.text = strings.Join(record, c.outputDelimiter)

	return true
}
