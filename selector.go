package main

import (
	"fmt"
)

type Selector interface {
	RequireHeader() bool
	ParseHeader(header []string) error
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

func (i *Indexes) RequireHeader() bool {
	return false
}

func (i *Indexes) ParseHeader(header []string) error {
	return nil
}

func (i *Indexes) Select(recode []string) ([]string, error) {
	var values []string
	for _, index := range i.indexes {
		values = append(values, recode[index])
	}
	return values, nil
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

func (h *Headers) RequireHeader() bool {
	return true
}

func (h *Headers) ParseHeader(headers []string) error {
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
	var values []string
	for _, index := range h.indexes {
		values = append(values, recode[index])
	}
	return values, nil
}
