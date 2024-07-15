package csv

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseCsvHeader(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected CsvHeader
	}{
		{
			name: "Simple case",
			line: "header1,header2,header3",
			expected: CsvHeader{
				headers: []string{"header1", "header2", "header3"},
			},
		},
		{
			name: "Extra spaces",
			line: " header1 , header2 , header3 ",
			expected: CsvHeader{
				headers: []string{" header1 ", " header2 ", " header3 "},
			},
		},
		{
			name: "Empty headers",
			line: ",,",
			expected: CsvHeader{
				headers: []string{"", "", ""},
			},
		},
		{
			name: "Single header",
			line: "header1",
			expected: CsvHeader{
				headers: []string{"header1"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseHeader(tt.line)
			assert.Equal(t, tt.expected, result)
		})
	}
}
