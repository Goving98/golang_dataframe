package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    // "strconv"
    // "time"
    "github.com/go-gota/gota/dataframe"
    // "github.com/go-gota/gota/series"
)

// Updated struct definitions for new JSON format
type TableData struct {
    Table struct {
        GroupVal string     `json:"groupVal"`
        Data     []TableRow `json:"data"`
    } `json:"table"`
}

type TableRow struct {
    ID      string `json:"id"`
    PageNo  int    `json:"page_no"`
    Type    string `json:"type"`
    Label   string `json:"label"`
    Cells   []Cell `json:"cells"`
}

type Cell struct {
    ID    string  `json:"id"`
    Label string  `json:"label"`
    Row   int     `json:"row"`
    Col   int     `json:"col"`
    XMin  float64 `json:"xmin"`
    XMax  float64 `json:"xmax"`
    YMin  float64 `json:"ymin"`
    YMax  float64 `json:"ymax"`
    Text  string  `json:"text"`
    Score float64 `json:"score"`
}

func main() {
    // Load and parse JSON
    jsonBytes, err := ioutil.ReadFile("getData.json")
    if err != nil {
        log.Fatalf("Error reading JSON file: %v", err)
    }

    var tableData TableData
    if err := json.Unmarshal(jsonBytes, &tableData); err != nil {
        log.Fatalf("Error decoding JSON: %v", err)
    }

    // Extract cells from the first table (assuming single table)
    if len(tableData.Table.Data) == 0 {
        log.Fatal("No table data found")
    }
    cells := tableData.Table.Data[0].Cells

    // Create maps for unique rows and columns
    rowMap := make(map[int]bool)
    colMap := make(map[int]bool)
    labelMap := make(map[int]string)

    // Get unique rows, columns and labels
    for _, cell := range cells {
        rowMap[cell.Row] = true
        colMap[cell.Col] = true
        labelMap[cell.Col] = cell.Label
    }

    // Create records for DataFrame
    records := make([][]string, 0)
    maxRow := len(rowMap)
    maxCol := len(colMap)

    // Create records matrix
    for row := 1; row <= maxRow; row++ {
        rowData := make([]string, maxCol+1) // +1 for ColumnOrder
        for _, cell := range cells {
            if cell.Row == row {
                rowData[cell.Col-1] = cell.Text
            }
        }
        // Add ColumnOrder
        rowData[maxCol] = fmt.Sprintf("col_%d", row)
        records = append(records, rowData)
    }

    // Create headers
    headers := make([]string, 0)
    for i := 1; i <= maxCol; i++ {
        if label, ok := labelMap[i]; ok {
            headers = append(headers, label)
        }
    }
    headers = append(headers, "ColumnOrder")

    // Create DataFrame
    df := dataframe.LoadRecords(records, dataframe.HasHeader(false))
    df.SetNames(headers...)

    // Display original data
    fmt.Println("Original DataFrame:")
    fmt.Println(df)

    // Rest of your operations (sort, filter, etc.) remain the same
    // ...existing code...
}