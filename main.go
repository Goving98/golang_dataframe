package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
)

// Structs for parsing JSON
type TableData struct {
	Table []Table `json:"table_1"`
}

type Table struct {
	Cells []Cell `json:"cells"`
}

type Cell struct {
	Row     int    `json:"row"`
	Col     int    `json:"col"`
	OCRText string `json:"ocr_text"`
}

func main() {
	// Sample JSON input
	jsonData := `
    {
    "type": "table",
    "table_1": [
            {
                "id": "d40bb083-0fc2-4e5c-8a11-ca4939e22b1a",
                "vertices":{
                    "xmin": 1168,
                    "xmax": 1381,
                    "ymin": 689,
                    "ymax": 724
                },
                "page_no": 0,
                "cells": [
                    {
                        "id": "3bedf90a-1e6c-4896-a28c-721dd7b753ad",
                        "label": "s_no",
                        "row": 5,
                        "col": 1,
                        "vertices":{
                            "xmin": 1168,
                            "xmax": 1381,
                            "ymin": 689,
                            "ymax": 724
                        },
                        "ocr_text": "1 ",
                        "score": 0
                    } , 
                    {
                        "id": "3bedf90a-1e6c-4896-a28c-721dd7b753ad",
                        "label": "s_no",
                        "row": 3,
                        "col": 2,
                        "vertices":{
                            "xmin": 1168,
                            "xmax": 1381,
                            "ymin": 689,
                            "ymax": 724
                        },
                        "ocr_text": "60.1",
                        "score": 0
                    }
                ]
            }
        ]
    }`

	// Parse JSON
	var tableData TableData
	err := json.Unmarshal([]byte(jsonData), &tableData)
	if err != nil {
		log.Fatal("Error parsing JSON:", err)
	}

	// Lists to store extracted data
	var ocrTexts []string
	var rowIds []int
	var colIds []int

	// Extract values from JSON
	for _, table := range tableData.Table {
		for _, cell := range table.Cells {
			ocrTexts = append(ocrTexts, cell.OCRText)
			rowIds = append(rowIds, cell.Row)
			colIds = append(colIds, cell.Col)
		}
	}

	// Create Gota DataFrame
	df := dataframe.New(
		series.New(ocrTexts, series.String, "OCR_TEXT"),
		series.New(rowIds, series.Int, "ROW_ID"),
		series.New(colIds, series.Int, "COL_ID"),
	)

	// Print the DataFrame
	fmt.Println(df)
    sortedDf := df.Arrange(dataframe.Sort("ROW_ID"))
    fmt.Println("\n📌 DataFrame after sorting by ROW_ID:")
    fmt.Println(sortedDf)
}
