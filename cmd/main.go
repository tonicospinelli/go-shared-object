package main

import (
	"fmt"
	"milenio.capital/code-challenge/pkg"
)

func main() {
	rawCsvData := "col1,col2,col3,col4,col5,col6,col7\nl1c1,l1c2,l1c3,l1c4,l1c5,l1c6,l1c7\nl1c1,l1c2,l1c3,l1c4,l1c5,l1c6,l1c7\nl2c1,l2c2,l2c3,l2c4,l2c5,l2c6,l2c7\nl3c1,l3c2,l3c3,l3c4,l3c5,l3c6,l3c7\n"
	showColumns := "col1,col3,col4,col7"
	filters := "col1>l1c1\ncol3>l1c3"

	fmt.Println("processCsv output:")
	pkg.ProcessCsv(rawCsvData, showColumns, filters)
	fmt.Println("")

	fmt.Println("processCsvFile output:")
	pkg.ProcessCsvFile("data.csv", showColumns, filters)
}
