package main

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
