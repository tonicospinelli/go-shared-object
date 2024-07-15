package csv

import "strings"

type CsvHeader struct {
	headers            []string
	selectedIndices    []int
	numSelectedColumns int
}

func (h *CsvHeader) Contains(s string) bool {
	for _, col := range h.headers {
		if col == s {
			return true
		}
	}
	return false
}

func parseHeader(line string) CsvHeader {
	headers := strings.Split(line, ",")
	return CsvHeader{
		headers: headers,
	}
}
