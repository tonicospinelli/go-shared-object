package csv

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseSelectedColumns(t *testing.T) {
	tests := []struct {
		name            string
		selectedColumns string
		csvHeader       CsvHeader
		expected        CsvHeader
	}{
		{
			name:            "All columns",
			selectedColumns: "",
			csvHeader:       CsvHeader{headers: []string{"header1", "header2", "header3"}},
			expected:        CsvHeader{headers: []string{"header1", "header2", "header3"}, selectedIndices: []int{0, 1, 2}, numSelectedColumns: 3},
		},
		{
			name:            "Select specific columns",
			selectedColumns: "header1,header3",
			csvHeader:       CsvHeader{headers: []string{"header1", "header2", "header3"}},
			expected:        CsvHeader{headers: []string{"header1", "header2", "header3"}, selectedIndices: []int{0, 2}, numSelectedColumns: 2},
		},
		{
			name:            "Select non-consecutive columns",
			selectedColumns: "header3,header1",
			csvHeader:       CsvHeader{headers: []string{"header1", "header2", "header3"}},
			expected:        CsvHeader{headers: []string{"header1", "header2", "header3"}, selectedIndices: []int{2, 0}, numSelectedColumns: 2},
		},
		{
			name:            "Select single column",
			selectedColumns: "header2",
			csvHeader:       CsvHeader{headers: []string{"header1", "header2", "header3"}},
			expected:        CsvHeader{headers: []string{"header1", "header2", "header3"}, selectedIndices: []int{1}, numSelectedColumns: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := parseSelectedColumns(tt.selectedColumns, &tt.csvHeader)
			assert.Nil(t, err)
			assert.Equal(t, tt.expected, tt.csvHeader)
		})
	}
}

func TestParseSelectedUnknownColumns(t *testing.T) {
	tests := []struct {
		name            string
		selectedColumns string
		csvHeader       CsvHeader
		expected        error
	}{
		{
			name:            "unknown single column",
			selectedColumns: "header0",
			csvHeader:       CsvHeader{headers: []string{"header1", "header2", "header3"}},
			expected:        fmt.Errorf("Header 'header0' not found in CSV file/string"),
		},
		{
			name:            "Select non-existent column",
			selectedColumns: "header1,header4",
			csvHeader:       CsvHeader{headers: []string{"header1", "header2", "header3"}},
			expected:        fmt.Errorf("Header 'header4' not found in CSV file/string"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := parseSelectedColumns(tt.selectedColumns, &tt.csvHeader)
			assert.Equal(t, tt.expected, err)
		})
	}
}

func TestSelectColumns(t *testing.T) {
	tests := []struct {
		name     string
		row      []string
		indices  []int
		expected []string
	}{
		{
			name:     "Select all columns",
			row:      []string{"1", "2", "3"},
			indices:  []int{0, 1, 2},
			expected: []string{"1", "2", "3"},
		},
		{
			name:     "Select specific columns",
			row:      []string{"1", "2", "3"},
			indices:  []int{0, 2},
			expected: []string{"1", "3"},
		},
		{
			name:     "Select single column",
			row:      []string{"1", "2", "3"},
			indices:  []int{1},
			expected: []string{"2"},
		},
		{
			name:     "Select non-consecutive columns",
			row:      []string{"1", "2", "3"},
			indices:  []int{2, 0},
			expected: []string{"3", "1"},
		},
		{
			name:     "Empty row",
			row:      []string{},
			indices:  []int{},
			expected: nil,
		},
		{
			name:     "Out of range index",
			row:      []string{"1", "2", "3"},
			indices:  []int{0, 3},
			expected: []string{"1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := selectColumns(tt.row, tt.indices)
			assert.Equal(t, tt.expected, result)
		})
	}
}
