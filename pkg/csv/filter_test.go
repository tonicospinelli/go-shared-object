package csv

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseFilters(t *testing.T) {
	tests := []struct {
		name                 string
		header               CsvHeader
		rowFilterDefinitions string
		expected             []Filter
	}{
		{
			name:                 "Single filter",
			header:               CsvHeader{headers: []string{"header1"}},
			rowFilterDefinitions: "header1=1",
			expected:             []Filter{{column: "header1", comparator: '=', value: "1"}},
		},
		{
			name:                 "Multiple filters",
			header:               CsvHeader{headers: []string{"header1", "header2", "header3"}},
			rowFilterDefinitions: "header1=1\nheader2>2\nheader3<3",
			expected:             []Filter{{column: "header1", comparator: '=', value: "1"}, {column: "header2", comparator: '>', value: "2"}, {column: "header3", comparator: '<', value: "3"}},
		},
		{
			name:                 "Empty filter definitions",
			header:               CsvHeader{headers: []string{"header1", "header2", "header3"}},
			rowFilterDefinitions: "",
			expected:             []Filter(nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseFilters(tt.rowFilterDefinitions, tt.header)
			assert.Nil(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseUnknownFilters(t *testing.T) {
	tests := []struct {
		name                 string
		header               CsvHeader
		rowFilterDefinitions string
		expected             error
	}{
		{
			name:                 "Invalid filter format",
			header:               CsvHeader{headers: []string{"header1", "header2", "header3"}},
			rowFilterDefinitions: "header1=1\nheader2>2\ninvalidfilter",
			expected:             fmt.Errorf("Header 'invalidfilter' not found in CSV file/string"),
		},
		{
			name:                 "Invalid column name",
			header:               CsvHeader{headers: []string{"header1", "header2", "header3"}},
			rowFilterDefinitions: "header1=1\nheader2>2\nheader4>0",
			expected:             fmt.Errorf("Header 'header4' not found in CSV file/string"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseFilters(tt.rowFilterDefinitions, tt.header)
			assert.Equal(t, tt.expected, err)
		})
	}
}

func TestApplyFilters(t *testing.T) {
	tests := []struct {
		name      string
		row       []string
		filters   []Filter
		csvHeader CsvHeader
		expected  bool
	}{
		{
			name:      "All filters match",
			row:       []string{"1", "2", "3"},
			filters:   []Filter{{column: "header1", comparator: '=', value: "1"}, {column: "header2", comparator: '=', value: "2"}},
			csvHeader: CsvHeader{headers: []string{"header1", "header2", "header3"}},
			expected:  true,
		},
		{
			name:      "Some filters do not match",
			row:       []string{"1", "2", "3"},
			filters:   []Filter{{column: "header1", comparator: '=', value: "1"}, {column: "header2", comparator: '=', value: "3"}},
			csvHeader: CsvHeader{headers: []string{"header1", "header2", "header3"}},
			expected:  false,
		},
		{
			name:      "No filters",
			row:       []string{"1", "2", "3"},
			filters:   []Filter{},
			csvHeader: CsvHeader{headers: []string{"header1", "header2", "header3"}},
			expected:  true,
		},
		{
			name:      "Filter on non-existent column",
			row:       []string{"1", "2", "3"},
			filters:   []Filter{{column: "header4", comparator: '=', value: "4"}},
			csvHeader: CsvHeader{headers: []string{"header1", "header2", "header3"}},
			expected:  false,
		},
		{
			name:      "Mixed filters with different comparators",
			row:       []string{"1", "2", "3"},
			filters:   []Filter{{column: "header1", comparator: '=', value: "1"}, {column: "header2", comparator: '>', value: "1"}, {column: "header3", comparator: '<', value: "4"}},
			csvHeader: CsvHeader{headers: []string{"header1", "header2", "header3"}},
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := applyFilters(tt.row, tt.filters, tt.csvHeader)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestApplyFilter(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		filter   Filter
		expected bool
	}{
		{
			name:     "Equal comparison - match",
			value:    "5",
			filter:   Filter{column: "header1", comparator: '=', value: "5"},
			expected: true,
		},
		{
			name:     "Equal comparison - no match",
			value:    "5",
			filter:   Filter{column: "header1", comparator: '=', value: "6"},
			expected: false,
		},
		{
			name:     "Greater than comparison - match",
			value:    "7",
			filter:   Filter{column: "header1", comparator: '>', value: "5"},
			expected: true,
		},
		{
			name:     "Greater than comparison - no match",
			value:    "5",
			filter:   Filter{column: "header1", comparator: '>', value: "5"},
			expected: false,
		},
		{
			name:     "Less than comparison - match",
			value:    "3",
			filter:   Filter{column: "header1", comparator: '<', value: "5"},
			expected: true,
		},
		{
			name:     "Less than comparison - no match",
			value:    "5",
			filter:   Filter{column: "header1", comparator: '<', value: "5"},
			expected: false,
		},
		{
			name:     "Invalid comparator",
			value:    "5",
			filter:   Filter{column: "header1", comparator: '#', value: "5"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := applyFilter(tt.value, tt.filter)
			if result != tt.expected {
				t.Errorf("applyFilter(%q, %v) = %v; want %v", tt.value, tt.filter, result, tt.expected)
			}
		})
	}
}
