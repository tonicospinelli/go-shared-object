package pkg

import "C"
import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"sync"
)

var mtx sync.Mutex

// DEPRECATED: use csv.ProcessCsv
func ProcessCsv(csvData string, selectedColumns string, rowFilterDefinitions string) {
	csvHeader := parseHeader(strings.Split(csvData, "\n")[0])
	parseSelectedColumns(selectedColumns, &csvHeader)
	filters := parseFilters(rowFilterDefinitions)
	processCsvData(csvData, csvHeader, filters)
}

// DEPRECATED: use csv.ProcessCsvFile
func ProcessCsvFile(csvFilePath string, selectedColumns string, rowFilterDefinitions string) {
	file, err := os.Open(csvFilePath)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to open file: %v\n", err)
	}
	defer func() {
		if file != nil {
			_ = file.Close()
		}
	}()

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to read file: %v\n", err)
	}

	csvData := strings.Join(flatten(lines), "\n")
	csvHeader := parseHeader(strings.Split(csvData, "\n")[0])
	parseSelectedColumns(selectedColumns, &csvHeader)
	filters := parseFilters(rowFilterDefinitions)
	processCsvData(csvData, csvHeader, filters)
}

func flatten(records [][]string) []string {
	var flat []string
	for _, record := range records {
		flat = append(flat, strings.Join(record, ","))
	}
	return flat
}

func processCsvData(csvData string, csvHeader CsvHeader, filters []Filter) {
	lines := strings.Split(csvData, "\n")
	mtx.Lock()
	defer mtx.Unlock()

	// Append headers
	fmt.Println(strings.Join(selectColumns(csvHeader.headers, csvHeader.selectedIndexes), ","))

	// Process each row
	for _, line := range lines[1:] {
		row := strings.Split(line, ",")
		if applyFilters(row, filters, csvHeader) {
			fmt.Println(strings.Join(selectColumns(row, csvHeader.selectedIndexes), ","))
		}
	}

}

type CsvHeader struct {
	headers            []string
	selectedIndexes    []int
	numSelectedColumns int
}

func parseHeader(line string) CsvHeader {
	headers := strings.Split(line, ",")
	return CsvHeader{
		headers: headers,
	}
}

type Filter struct {
	column     string
	comparator byte
	value      string
}

func parseFilters(rowFilterDefinitions string) []Filter {
	var filters []Filter
	lines := strings.Split(rowFilterDefinitions, "\n")
	for _, line := range lines {
		for i, char := range line {
			if char == '=' || char == '>' || char == '<' {
				filters = append(filters, Filter{
					column:     line[:i],
					comparator: byte(char),
					value:      line[i+1:],
				})
				break
			}
		}
	}
	return filters
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

func parseSelectedColumns(selectedColumns string, csvHeader *CsvHeader) {
	columns := strings.Split(selectedColumns, ",")
	if columns[0] == "" {
		columns = csvHeader.headers
	}
	for _, col := range columns {
		for i, header := range csvHeader.headers {
			if header == col {
				csvHeader.selectedIndexes = append(csvHeader.selectedIndexes, i)
				break
			}
		}
	}
	csvHeader.numSelectedColumns = len(csvHeader.selectedIndexes)
}

func selectColumns(row []string, indexes []int) []string {
	var selected []string
	length := len(row)
	for _, idx := range indexes {
		if idx < length {
			selected = append(selected, row[idx])
		}
	}
	return selected
}
