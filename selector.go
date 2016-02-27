package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Selector interface {
	DropHeaders() bool
	ParseHeaders(headers []string) error
	Select(record []string) ([]string, error)
}

type All struct {
}

func NewAll() *All {
	return &All{}
}

func (a *All) DropHeaders() bool {
	return false
}

func (a *All) ParseHeaders(headers []string) error {
	return nil
}

func (a *All) Select(record []string) ([]string, error) {
	return record, nil
}

var (
	exprIndexes = regexp.MustCompile(`^(?:\d*-\d*|\d+)(?:,(?:\d*-\d*|\d+))*$`)
	exprIndex   = regexp.MustCompile(`(?:\d*-\d*|\d+)`)
	exprRange   = regexp.MustCompile(`^(\d*)-(\d*)$`)
)

func toIndex(s string) (index int, err error) {
	index, err = strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	if index == 0 {
		return 0, fmt.Errorf("indexes are numberd from 1")
	}
	return index, nil
}

type Indexes struct {
	list    string
	indexes []int
}

func NewIndexes(list string) *Indexes {
	return &Indexes{
		list: list,
	}
}

func (i *Indexes) DropHeaders() bool {
	return false
}

func (i *Indexes) ParseHeaders(headers []string) error {
	if i.list == "" {
		i.indexes = make([]int, 0)
		return nil
	}
	if !exprIndexes.MatchString(i.list) {
		return fmt.Errorf("%q: invalid syntax", i.list)
	}

	i.indexes = make([]int, 0)
	for _, rawIndex := range exprIndex.FindAllString(i.list, -1) {
		var err error
		switch {
		case exprRange.MatchString(rawIndex):
			first, last := 1, len(headers)
			rawRange := exprRange.FindStringSubmatch(rawIndex)
			if rawRange[1] != "" {
				first, err = toIndex(rawRange[1])
				if err != nil {
					return err
				}
			}
			if rawRange[2] != "" {
				last, err = toIndex(rawRange[2])
				if err != nil {
					return err
				}
			}
			for index := first; index <= last && index <= len(headers); index++ {
				i.indexes = append(i.indexes, index-1)
			}
		default:
			index, err := toIndex(rawIndex)
			if err != nil {
				return err
			}
			i.indexes = append(i.indexes, index-1)
		}
	}
	return nil
}

func (i *Indexes) Select(record []string) ([]string, error) {
	a := make([]string, len(i.indexes))
	for j, index := range i.indexes {
		if index >= 0 && index < len(record) {
			a[j] = record[index]
		}
	}
	return a, nil
}

var (
	exprHeader    = regexp.MustCompile(`(?:[^,\\]|\\.)*`)
	exprBackslash = regexp.MustCompile(`\\(.)`)
	exprTrailing  = regexp.MustCompile(`\\+$`)
)

type Headers struct {
	indexes []int
	headers []string
}

func NewHeaders(list string) *Headers {
	list = exprTrailing.ReplaceAllStringFunc(list, func(s string) string {
		return strings.Repeat(`\\`, len(s)/2)
	})
	if list == "" {
		return &Headers{
			headers: []string{},
		}
	}

	headers := exprHeader.FindAllString(list, -1)
	for i := 0; i < len(headers); i++ {
		headers[i] = exprBackslash.ReplaceAllString(headers[i], "$1")
	}
	return &Headers{
		headers: headers,
	}
}

func (h *Headers) DropHeaders() bool {
	return true
}

func (h *Headers) ParseHeaders(headers []string) error {
	indexMap := make(map[string]int)
	for i, header := range headers {
		if _, ok := indexMap[header]; ok {
			return fmt.Errorf("%q: duplicated header", header)
		}
		indexMap[header] = i
	}

	h.indexes = make([]int, len(h.headers))
	for i, header := range h.headers {
		if index, ok := indexMap[header]; ok {
			h.indexes[i] = index
		} else {
			h.indexes[i] = -1
		}
	}
	return nil
}

func (h *Headers) Select(record []string) ([]string, error) {
	a := make([]string, len(h.indexes))
	for i, index := range h.indexes {
		if index != -1 {
			a[i] = record[index]
		}
	}
	return a, nil
}
