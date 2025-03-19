package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"table_transform/helpers"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
)

// Structs for transformed JSON
type TableData struct {
	Type     string            `json:"type"`
	Metadata Metadata          `json:"metadata"`
	Cells    []Cell            `json:"cells"`
	Layout   Layout            `json:"layout"`
	Headers  map[string]string `json:"headers"`
}

type Metadata struct {
	ID       string   `json:"id"`
	PageNo   int      `json:"page_no"`
	Vertices Vertices `json:"vertices"`
}

type Vertices struct {
	XMin float64 `json:"xmin"`
	XMax float64 `json:"xmax"`
	YMin float64 `json:"ymin"`
	YMax float64 `json:"ymax"`
}

type Cell struct {
	ID         string      `json:"id"`
	RowID      string      `json:"row_id"`
	ColID      string      `json:"col_id"`
	Vertices   Vertices    `json:"vertices"`
	OcrText    interface{} `json:"ocr_text"`
	IsHeader   bool        `json:"is_header"`
	Confidence float64     `json:"confidence"`
}

type Layout struct {
	RowOrder    []string `json:"row_order"`
	ColumnOrder []string `json:"column_order"`
}


func main() {
	// Load transformed JSON
	jsonBytes, err := ioutil.ReadFile("getData.json")
	if err != nil {
		log.Fatalf("Error reading JSON file: %v", err)
	}

	var tableData TableData
	if err := json.Unmarshal(jsonBytes, &tableData); err != nil {
		log.Fatalf("Error decoding JSON: %v", err)
	}

	// Extract headers in order
	headers := make([]string, len(tableData.Layout.ColumnOrder))
	for i, colID := range tableData.Layout.ColumnOrder {
		headers[i] = tableData.Headers[colID]
	}

	// Extract row data
	records := make([][]interface{}, len(tableData.Layout.RowOrder))
	for i, rowID := range tableData.Layout.RowOrder {
		row := make([]interface{}, len(tableData.Layout.ColumnOrder))
		for j, colID := range tableData.Layout.ColumnOrder {
			for _, cell := range tableData.Cells {
				if cell.RowID == rowID && cell.ColID == colID {
					row[j] = cell.OcrText
				}
			}
		}
		records[i] = row
	}

	// Create DataFrame using series
	seriesList := make([]series.Series, len(headers))
	for j, header := range headers {
		colData := make([]interface{}, len(records))
		for i, row := range records {
			colData[i] = row[j]
		}
		seriesList[j] = series.New(colData, series.String, header)
	}

	df := dataframe.New(seriesList...)
	// var df dataframe.DataFrame
	fmt.Println("Before DataFrame:")
	fmt.Println(df)
	// Assume df is preloaded with data

	// df = HandleOperations(df)
	// fmt.Println("Final DataFrame:")
	// fmt.Println(df)

	db, ctx, err :=helpers.ConnectMongoDB()
	if err != nil {
		fmt.Println("Error connecting to MongoDB:", err)
		return
	}

	err = helpers.StoreTableDataAndConfig(db, ctx, df)
	if err != nil {
		fmt.Println("Error storing data:", err)
		return
	}

	fmt.Println("Updated DataFrame successfully stored in MongoDB.")
}
