package csv

import (
	"fmt"
	"strings"
)

type Filter struct {
	column     string
	comparator byte
	value      string
}

func NewFilter(column string, comparator byte, value string) *Filter {
	return &Filter{
		column:     column,
		comparator: comparator,
		value:      value,
	}
}

func ParseFilters(f string, h CsvHeader) ([]Filter, error) {
	if f == "" {
		return nil, nil
	}
	var filters []Filter
	lines := strings.Split(f, "\n")
	for _, line := range lines {
		filter, err := parseFilter(line, h)
		if err != nil {
			return nil, err
		}
		if filter != nil {
			filters = append(filters, *filter)
		}
	}
	return filters, nil
}

func parseFilter(line string, h CsvHeader) (*Filter, error) {
	for i, char := range line {
		if char == '=' || char == '>' || char == '<' {
			col := line[:i]
			val := line[i+1:]
			if h.Contains(col) {
				return NewFilter(col, byte(char), val), nil
			} else {
				return nil, fmt.Errorf("Header '%s' not found in CSV file/string", col)
			}
		}
	}
	return nil, fmt.Errorf("Header '%s' not found in CSV file/string", line)
}

func applyFilter(value string, filter Filter) bool {
	switch filter.comparator {
	case '=':
		return strings.Compare(value, filter.value) == 0
	case '>':
		return strings.Compare(value, filter.value) > 0
	case '<':
		return strings.Compare(value, filter.value) < 0
	default:
		return false
	}
}

func applyFilters(row []string, filters []Filter, csvHeader CsvHeader) bool {
	for _, filter := range filters {
		columnIndex := -1
		for i, header := range csvHeader.headers {
			if header == filter.column {
				columnIndex = i
				break
			}
		}
		if columnIndex == -1 || !applyFilter(row[columnIndex], filter) {
			return false
		}
	}
	return true
}
