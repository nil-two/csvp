package main

type Selector interface {
	RequireHeader() bool
	ParseHeader(header []string) error
	Select(recode []string) ([]string, error)
}
