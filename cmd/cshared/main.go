package main

import "C"
import (
	"fmt"
	"milenio.capital/code-challenge/pkg/csv"
	"os"
)

func main() {}

//export processCsv
func processCsv(csvData *C.char, selectedColumns *C.char, rowFilterDefinitions *C.char) {
	goCsv := C.GoString(csvData)
	goSelectedColumns := C.GoString(selectedColumns)
	goRowFilterDefinitions := C.GoString(rowFilterDefinitions)
	err := csv.ProcessCsv(goCsv, goSelectedColumns, goRowFilterDefinitions)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
}

//export processCsvFile
func processCsvFile(csvFilePath *C.char, selectedColumns *C.char, rowFilterDefinitions *C.char) {
	goCsvFilePath := C.GoString(csvFilePath)
	goSelectedColumns := C.GoString(selectedColumns)
	goRowFilterDefinitions := C.GoString(rowFilterDefinitions)
	err := csv.ProcessCsvFile(goCsvFilePath, goSelectedColumns, goRowFilterDefinitions)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
}
