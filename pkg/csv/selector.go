package csv

import (
	"fmt"
	"strings"
)

func parseSelectedColumns(selectedColumns string, csvHeader *CsvHeader) error {
	columns := strings.Split(strings.Trim(selectedColumns, ""), ",")

	if columns[0] == "" {
		columns = csvHeader.headers
	}

	for _, col := range columns {
		unknown := true
		for i, header := range csvHeader.headers {
			if header == col {
				csvHeader.selectedIndices = append(csvHeader.selectedIndices, i)
				unknown = false
				break
			}
		}
		if unknown {
			return fmt.Errorf("Header '%s' not found in CSV file/string", col)
		}
	}
	csvHeader.numSelectedColumns = len(csvHeader.selectedIndices)
	return nil
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
