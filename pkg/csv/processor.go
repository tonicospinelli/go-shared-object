package csv

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

func processCsvData(csvData string, csvHeader CsvHeader, filters []Filter) {
	lines := strings.Split(csvData, "\n")

	fmt.Println(strings.Join(selectColumns(csvHeader.headers, csvHeader.selectedIndices), ","))

	for _, line := range lines[1:] {
		row := strings.Split(line, ",")
		if applyFilters(row, filters, csvHeader) {
			fmt.Println(strings.Join(selectColumns(row, csvHeader.selectedIndices), ","))
		}
	}
}

func ProcessCsv(csvData string, selectedColumns string, rowFilterDefinitions string) error {
	csvHeader := parseHeader(strings.Split(csvData, "\n")[0])
	err := parseSelectedColumns(selectedColumns, &csvHeader)
	if err != nil {
		return err
	}
	filters, err := ParseFilters(rowFilterDefinitions, csvHeader)
	processCsvData(csvData, csvHeader, filters)
	return nil
}

func ProcessCsvFile(csvFilePath string, selectedColumns string, rowFilterDefinitions string) error {
	file, err := os.Open(csvFilePath)
	if err != nil {
		return fmt.Errorf("Failed to open file %s:, error:%v\n", csvFilePath, err)
	}
	defer func() { _ = file.Close() }()

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("Failed to read file %s: error %v\n", csvFilePath, err)
	}

	csvData := strings.Join(flatten(lines), "\n")
	return ProcessCsv(csvData, selectedColumns, rowFilterDefinitions)
}

func flatten(records [][]string) []string {
	var flat []string
	for _, record := range records {
		flat = append(flat, strings.Join(record, ","))
	}
	return flat
}
