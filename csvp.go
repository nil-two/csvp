package main

import (
	"regexp"
	"strconv"
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
