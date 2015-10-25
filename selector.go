package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var FIELD = regexp.MustCompile(`(?:[^,\\]|\\.)*`)

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

type Indexes struct {
	indexes []int
}

func NewIndexes(list string) (*Indexes, error) {
	fields := FIELD.FindAllString(list, -1)

	indexes := make([]int, len(fields))
	for i, field := range fields {
		index, err := strconv.Atoi(field)
		if err != nil {
			return nil, err
		}
		if index < 1 {
			return nil, fmt.Errorf("indexes are numberd from 1")
		}
		indexes[i] = index - 1
	}
	return &Indexes{
		indexes: indexes,
	}, nil
}

func (i *Indexes) DropHeaders() bool {
	return false
}

func (i *Indexes) ParseHeaders(headers []string) error {
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

type Headers struct {
	indexes []int
	headers []string
}

func NewHeaders(list string) (*Headers, error) {
	headers := FIELD.FindAllString(list, -1)
	for i := 0; i < len(headers); i++ {
		headers[i] = strings.Replace(headers[i], `\,`, `,`, -1)
	}
	return &Headers{
		headers: headers,
	}, nil
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
