package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Selector interface {
	RequireHeaders() bool
	ParseHeaders(headers []string) error
	Select(recode []string) ([]string, error)
}

type Indexes struct {
	indexes []int
}

func NewIndexes(indexes []int) *Indexes {
	return &Indexes{
		indexes: indexes,
	}
}

func (i *Indexes) RequireHeaders() bool {
	return false
}

func (i *Indexes) ParseHeaders(headers []string) error {
	return nil
}

func (i *Indexes) Select(recode []string) ([]string, error) {
	a := make([]string, len(i.indexes))
	for j, index := range i.indexes {
		if index >= 0 || index < len(recode) {
			a[j] = recode[index]
		}
	}
	return a, nil
}

type Headers struct {
	indexes []int
	headers []string
}

func NewHeaders(headers []string) *Headers {
	return &Headers{
		headers: headers,
	}
}

func (h *Headers) RequireHeaders() bool {
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

var FIELD = regexp.MustCompile(`(?:[^,\\]|\\.)*`)

func parseIndexesList(list string) ([]int, error) {
	fields := FIELD.FindAllString(list, -1)
	indexes := make([]int, len(fields))
	for i, field := range fields {
		index, err := strconv.Atoi(field)
		if err != nil {
			return nil, err
		}
		indexes[i] = index - 1
	}
	return indexes, nil
}

func parseHeadersList(list string) ([]string, error) {
	fields := FIELD.FindAllString(list, -1)
	for i := 0; i < len(fields); i++ {
		fields[i] = strings.Replace(fields[i], `\,`, `,`, -1)
	}
	return fields, nil
}
