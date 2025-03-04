package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
    // "go_learn/cmd/transformer/transform"
	"log"
	"github.com/go-gota/gota/dataframe"
)

type TransformedTable struct {
	Cells []struct {
		RowID   string `json:"row_id"`
		ColID   string `json:"col_id"`
		OcrText string `json:"ocr_text"`
	} `json:"cells"`
}

func main() {
	// Read transformed JSON
	jsonBytes, err := ioutil.ReadFile("data/getData.json")
	if err != nil {
		log.Fatalf("Error reading transformed JSON: %v", err)
	}

	// Parse JSON
	var tableData TransformedTable
	err = json.Unmarshal(jsonBytes, &tableData)
	if err != nil {
		log.Fatalf("Error decoding JSON: %v", err)
	}

	// Build DataFrame
	records := make([][]string, 0)
	headers := map[string]bool{}
	for _, cell := range tableData.Cells {
		headers[cell.ColID] = true
	}

	// Collect unique column headers
	headerRow := make([]string, 0)
	for col := range headers {
		headerRow = append(headerRow, col)
	}

	// Fill rows
	rowMap := make(map[string]map[string]string)
	for _, cell := range tableData.Cells {
		if _, exists := rowMap[cell.RowID]; !exists {
			rowMap[cell.RowID] = make(map[string]string)
		}
		rowMap[cell.RowID][cell.ColID] = cell.OcrText
	}

	// Convert to slice format
	for _, rowID := range rowMap {
		row := make([]string, len(headerRow))
		for i, col := range headerRow {
			row[i] = rowID[col]
		}
		records = append(records, row)
	}

	// Create DataFrame
	df := dataframe.LoadRecords(records, dataframe.HasHeader(false))
	df.SetNames(headerRow...)

	fmt.Println("✅ DataFrame Loaded:")
	fmt.Println(df)
}
