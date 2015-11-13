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
	Select(recode []string) ([]string, error)
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

func (a *All) Select(recode []string) ([]string, error) {
	return recode, nil
}

var (
	INDEXES = regexp.MustCompile(`^(?:\d*-\d*|\d+)(?:,(?:\d*-\d*|\d+))*$`)
	INDEX   = regexp.MustCompile(`(?:\d*-\d*|\d+)`)
	RANGE   = regexp.MustCompile(`^(\d*)-(\d*)$`)
)

func toIndex(s string) (index int, err error) {
	index, err = strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	if index == 0 {
		return 0, fmt.Errorf("indexes are numberd from 1")
	}
	return index, err
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
	if !INDEXES.MatchString(i.list) {
		return fmt.Errorf("%q: invalid syntax", i.list)
	}

	i.indexes = make([]int, 0)
	for _, rawIndex := range INDEX.FindAllString(i.list, -1) {
		var err error
		switch {
		case RANGE.MatchString(rawIndex):
			first, last := 1, len(headers)
			rawRange := RANGE.FindStringSubmatch(rawIndex)
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

func (i *Indexes) Select(recode []string) ([]string, error) {
	a := make([]string, len(i.indexes))
	for j, index := range i.indexes {
		if index >= 0 && index < len(recode) {
			a[j] = recode[index]
		}
	}
	return a, nil
}

var HEADER = regexp.MustCompile(`(?:[^,\\]|\\.)*`)

type Headers struct {
	indexes []int
	headers []string
}

func NewHeaders(list string) *Headers {
	if list == "" {
		return &Headers{
			headers: []string{},
		}
	}

	headers := HEADER.FindAllString(list, -1)
	for i := 0; i < len(headers); i++ {
		headers[i] = strings.Replace(headers[i], `\,`, `,`, -1)
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

func (h *Headers) Select(recode []string) ([]string, error) {
	a := make([]string, len(h.indexes))
	for i, index := range h.indexes {
		if index != -1 {
			a[i] = recode[index]
		}
	}
	return a, nil
}
