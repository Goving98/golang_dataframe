package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "strconv"
	"time"
    "github.com/go-gota/gota/dataframe"
    "github.com/go-gota/gota/series"
)

// Updated struct definitions for new JSON format
type TableData struct {
    Type     string            `json:"type"`
    Metadata struct {
        ID       string   `json:"id"`
        PageNo   int      `json:"page_no"`
        Vertices Vertices `json:"vertices"`
    } `json:"metadata"`
    Cells   []Cell            `json:"cells"`
    Layout  Layout            `json:"layout"`
    Headers map[string]string `json:"headers"`
}

type Vertices struct {
    XMin int `json:"xmin"`
    XMax int `json:"xmax"`
    YMin int `json:"ymin"`
    YMax int `json:"ymax"`
}

type Cell struct {
    ID         string   `json:"id"`
    RowID      string   `json:"row"`
    ColID      string   `json:"col"`
    Vertices   Vertices `json:"vertices"`
    OcrText    string   `json:"text"`
    IsHeader   bool     `json:"is_header"`
    Confidence float64  `json:"confidence"`
}

type Layout struct {
    RowOrder    []string `json:"row_order"`
    ColumnOrder []string `json:"column_order"`
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

    // Create DataFrame from JSON data
    headers := make([]string, 0, len(tableData.Headers))
    for _, colID := range tableData.Layout.ColumnOrder {
        headers = append(headers, tableData.Headers[colID])
    }

    records := make([][]string, 0)
    for _, rowID := range tableData.Layout.RowOrder {
        if rowID != "row_1" {
            row := make([]string, 0, len(tableData.Layout.ColumnOrder))
            for _, colID := range tableData.Layout.ColumnOrder {
                for _, cell := range tableData.Cells {
                    if cell.RowID == rowID && cell.ColID == colID && !cell.IsHeader {
                        row = append(row, cell.OcrText)
                    }
                }
            }
            // colID := 0
			i:=0
            row = append(row, strconv.Itoa(i))
            records = append(records, row)
        }
    }

    df := dataframe.LoadRecords(records, dataframe.HasHeader(false))
	headers = append(headers, "ColumnOrder")
    df.SetNames(headers...)

    // Create DataFrame using series
    df := dataframe.New(seriesList...)

    // Print arrays and DataFrame
    fmt.Println("Row Order:", rowOrder)
    fmt.Println("Column Order:", columnOrder)
    fmt.Println("\nDataFrame:")
    fmt.Println(df)

    // Sort by Price (descending)
	start := time.Now()
    sortedDF := df.Arrange(dataframe.RevSort("Price"))
	elapsed := time.Since(start)
    fmt.Println("\nSorted by Price (Descending):")
	fmt.Printf("Sorting time: %s\n", elapsed)
    fmt.Println(sortedDF)

	start = time.Now()
    sortedDF2 := df.Arrange(dataframe.Sort("Price"))
	elapsed = time.Since(start)
    fmt.Println("\nSorted by Price (Ascending):")
	fmt.Printf("Sorting time: %s\n", elapsed)
    fmt.Println(sortedDF2)

    // Filter products with price > 900
    filter := dataframe.F{Colname: "Price", Comparator: series.Greater, Comparando: 900.0}
	start = time.Now()
    filteredDF := df.Filter(filter)
    fmt.Println("\nFiltered (Price > 900):")
	elapsed = time.Since(start)
	fmt.Printf("Filtering time: %s\n", elapsed)
    fmt.Println(filteredDF)

	start = time.Now()
	searchDF := df.Filter(dataframe.F{Colname: "Product", Comparator: series.Eq, Comparando: "Laptop"})
	elapsed = time.Since(start)
	fmt.Printf("Search time: %s\n", elapsed)
    fmt.Println("\nSearch for ProductA:")
    fmt.Println(searchDF)

    // Insert a new row
    newRow := dataframe.New(
        series.New([]string{"NewProduct"}, series.String, "Product"),
        series.New([]float64{1200.0}, series.Float, "Price"),
		series.New([]string{"col_" + strconv.Itoa(len(tableData.Layout.ColumnOrder)+1)}, series.String, "ColumnOrder"),
    )
	start = time.Now()
    df = df.RBind(newRow)
	elapsed = time.Since(start)
	tableData.Layout.ColumnOrder = append(tableData.Layout.ColumnOrder, "col_"+strconv.Itoa(len(tableData.Layout.ColumnOrder)+1))
	fmt.Printf("Inserting time: %s\n", elapsed)
    fmt.Println("\nAfter Inserting New Row:")
    fmt.Println(df)

    // Split the DataFrame into two parts
    nRows := df.Nrow()
    middleIndex := nRows / 2
    
    // Create two DataFrames by filtering rows
	start = time.Now()
    df1 := df.Subset([]int{0, middleIndex - 1})
    df2 := df.Subset([]int{middleIndex, nRows - 1})
	elapsed = time.Since(start)
	fmt.Printf("Splitting time: %s\n", elapsed)
    fmt.Println("\nFirst Part of Split DataFrame:")
    fmt.Println(df1)
    fmt.Println("\nSecond Part of Split DataFrame:")
    fmt.Println(df2)

    // Join two DataFrames
	start =time.Now()
    // strings := 0
    joinedDF := df1.InnerJoin(df2,"Product")
	elapsed = time.Since(start)
	fmt.Printf("Joining time: %s\n", elapsed)
    fmt.Println("\nJoined DataFrame:")
    fmt.Println(joinedDF)
	
}