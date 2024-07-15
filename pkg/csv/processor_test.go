package csv

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

func TestProcessCsvData(t *testing.T) {
	tests := []struct {
		name      string
		csvData   string
		csvHeader CsvHeader
		filters   []Filter
		expected  string
	}{
		{
			name:      "No filters, all columns",
			csvData:   "header1,header2,header3\n1,2,3\n4,5,6",
			csvHeader: CsvHeader{},
			filters:   []Filter{},
			expected:  "header1,header2,header3\n1,2,3\n4,5,6\n",
		},
		{
			name:      "With filters, all columns",
			csvData:   "header1,header2,header3\n1,2,3\n4,5,6",
			csvHeader: CsvHeader{},
			filters:   []Filter{{column: "header1", comparator: '>', value: "1"}},
			expected:  "header1,header2,header3\n4,5,6\n",
		},
		{
			name:      "With filters, selected columns",
			csvData:   "header1,header2,header3\n1,2,3\n4,5,6",
			csvHeader: CsvHeader{},
			filters:   []Filter{{column: "header1", comparator: '>', value: "1"}},
			expected:  "header1,header2,header3\n4,5,6\n",
		},
		{
			name:      "With filters, different columns",
			csvData:   "header1,header2,header3\n1,2,3\n4,5,6",
			csvHeader: CsvHeader{},
			filters:   []Filter{{column: "header2", comparator: '=', value: "5"}},
			expected:  "header1,header2,header3\n4,5,6\n",
		},
		{
			name:      "Empty CSV data",
			csvData:   "",
			csvHeader: CsvHeader{},
			filters:   []Filter{},
			expected:  "\n",
		},
		{
			name:      "No match filters",
			csvData:   "header1,header2,header3\n1,2,3\n4,5,6",
			csvHeader: CsvHeader{},
			filters:   []Filter{{column: "header1", comparator: '=', value: "10"}},
			expected:  "header1,header2,header3\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processCsvData(tt.csvData, tt.csvHeader, tt.filters)
		})
	}
}

func TestFlatten(t *testing.T) {
	tests := []struct {
		name     string
		records  [][]string
		expected []string
	}{
		{
			name: "Multiple records",
			records: [][]string{
				{"header1", "header2", "header3"},
				{"1", "2", "3"},
				{"4", "5", "6"},
			},
			expected: []string{
				"header1,header2,header3",
				"1,2,3",
				"4,5,6",
			},
		},
		{
			name: "Single record",
			records: [][]string{
				{"header1", "header2", "header3"},
			},
			expected: []string{
				"header1,header2,header3",
			},
		},
		{
			name:     "Empty records",
			records:  [][]string{},
			expected: nil,
		},
		{
			name: "Empty strings in records",
			records: [][]string{
				{"header1", "header2", "header3"},
				{"", "", ""},
				{"4", "5", "6"},
			},
			expected: []string{
				"header1,header2,header3",
				",,",
				"4,5,6",
			},
		},
		{
			name: "Records with spaces",
			records: [][]string{
				{" header1 ", " header2 ", " header3 "},
				{" 1 ", " 2 ", " 3 "},
			},
			expected: []string{
				" header1 , header2 , header3 ",
				" 1 , 2 , 3 ",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := flatten(tt.records)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func captureOutput(f func()) string {
	r, w, _ := os.Pipe()
	stdout := os.Stdout
	os.Stdout = w

	f()

	_ = w.Close()
	os.Stdout = stdout

	out, _ := io.ReadAll(r)
	return string(out)
}

func captureError(f func()) string {
	r, w, _ := os.Pipe()
	stderr := os.Stderr
	os.Stderr = w

	f()

	_ = w.Close()
	os.Stderr = stderr

	out, _ := io.ReadAll(r)
	return string(out)
}

func TestProcessCsv(t *testing.T) {
	tests := []struct {
		name                 string
		csvData              string
		selectedColumns      string
		rowFilterDefinitions string
		expected             string
	}{
		{
			name:                 "No filters, all columns",
			csvData:              "header1,header2,header3\n1,2,3\n4,5,6",
			selectedColumns:      "",
			rowFilterDefinitions: "",
			expected:             "header1,header2,header3\n1,2,3\n4,5,6\n",
		},
		{
			name:                 "With filters, all columns",
			csvData:              "header1,header2,header3\n1,2,3\n4,5,6",
			selectedColumns:      "",
			rowFilterDefinitions: "header1>1",
			expected:             "header1,header2,header3\n4,5,6\n",
		},
		{
			name:                 "With filters, selected columns",
			csvData:              "header1,header2,header3\n1,2,3\n4,5,6",
			selectedColumns:      "header1,header3",
			rowFilterDefinitions: "header1>1",
			expected:             "header1,header3\n4,6\n",
		},
		{
			name:                 "With filters, different columns",
			csvData:              "header1,header2,header3\n1,2,3\n4,5,6",
			selectedColumns:      "header2,header3",
			rowFilterDefinitions: "header1>1",
			expected:             "header2,header3\n5,6\n",
		},
		{
			name:                 "Empty CSV data",
			csvData:              "",
			selectedColumns:      "",
			rowFilterDefinitions: "",
			expected:             "\n",
		},
		{
			name:                 "No match filters",
			csvData:              "header1,header2,header3\n1,2,3\n4,5,6",
			selectedColumns:      "",
			rowFilterDefinitions: "header1=10",
			expected:             "header1,header2,header3\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := captureOutput(func() {
				err := ProcessCsv(tt.csvData, tt.selectedColumns, tt.rowFilterDefinitions)
				if err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "Failed to process file: %v\n", err)
				}
			})
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestProcessCsvError(t *testing.T) {
	tests := []struct {
		name                 string
		csvData              string
		selectedColumns      string
		rowFilterDefinitions string
		expected             string
	}{
		{
			name:                 "unknown single column",
			csvData:              "header1,header2,header3\n1,2,3\n4,5,6",
			selectedColumns:      "header0",
			rowFilterDefinitions: "",
			expected:             "Header 'header0' not found in CSV file/string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := captureError(func() {
				err := ProcessCsv(tt.csvData, tt.selectedColumns, tt.rowFilterDefinitions)
				if err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "%v", err)
				}
			})
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestProcessCsvFile(t *testing.T) {
	tests := []struct {
		name                 string
		fileContent          string
		selectedColumns      string
		rowFilterDefinitions string
		expectedOutput       string
	}{
		{
			name:                 "No filters, all columns",
			fileContent:          "header1,header2,header3\n1,2,3\n4,5,6\n",
			selectedColumns:      "",
			rowFilterDefinitions: "",
			expectedOutput:       "header1,header2,header3\n1,2,3\n4,5,6\n",
		},
		{
			name:                 "With filters, all columns",
			fileContent:          "header1,header2,header3\n1,2,3\n4,5,6\n",
			selectedColumns:      "",
			rowFilterDefinitions: "header1>1",
			expectedOutput:       "header1,header2,header3\n4,5,6\n",
		},
		{
			name:                 "With filters, selected columns",
			fileContent:          "header1,header2,header3\n1,2,3\n4,5,6\n",
			selectedColumns:      "header1,header3",
			rowFilterDefinitions: "header1>1",
			expectedOutput:       "header1,header3\n4,6\n",
		},
		{
			name:                 "With filters, different columns",
			fileContent:          "header1,header2,header3\n1,2,3\n4,5,6\n",
			selectedColumns:      "header2,header3",
			rowFilterDefinitions: "header1>1",
			expectedOutput:       "header2,header3\n5,6\n",
		},
		{
			name:                 "Empty CSV data",
			fileContent:          ``,
			selectedColumns:      "",
			rowFilterDefinitions: "",
			expectedOutput:       "\n",
		},
		{
			name:                 "No match filters",
			fileContent:          "header1,header2,header3\n1,2,3\n4,5,6\n",
			selectedColumns:      "",
			rowFilterDefinitions: "header1=10",
			expectedOutput:       "header1,header2,header3\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpfile, err := os.CreateTemp("", "test.csv")
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				_ = os.Remove(tmpfile.Name())
			}()

			if _, err := tmpfile.Write([]byte(tt.fileContent)); err != nil {
				t.Fatal(err)
			}
			if err := tmpfile.Close(); err != nil {
				t.Fatal(err)
			}

			output := captureOutput(func() {
				_ = ProcessCsvFile(tmpfile.Name(), tt.selectedColumns, tt.rowFilterDefinitions)
			})

			assert.Equal(t, tt.expectedOutput, output)
		})
	}
}
