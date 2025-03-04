package transformer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type InputTable struct {
	Table struct {
		GroupVal string `json:"groupVal"`
		Data     []struct {
			ID       string `json:"id"`
			PageNo   int    `json:"page_no"`
			Cells    []struct {
				ID        string  `json:"id"`
				Label     string  `json:"label"`
				Row       int     `json:"row"`
				Col       int     `json:"col"`
				Xmin      float64 `json:"xmin"`
				Xmax      float64 `json:"xmax"`
				Ymin      float64 `json:"ymin"`
				Ymax      float64 `json:"ymax"`
				Text      string  `json:"text"`
				Score     float64 `json:"score"`
			} `json:"cells"`
		} `json:"data"`
	} `json:"table"`
}

type TransformedTable struct {
	Type     string `json:"type"`
	Metadata struct {
		ID      string `json:"id"`
		PageNo  int    `json:"page_no"`
		Vertices struct {
			Xmin float64 `json:"xmin"`
			Xmax float64 `json:"xmax"`
			Ymin float64 `json:"ymin"`
			Ymax float64 `json:"ymax"`
		} `json:"vertices"`
	} `json:"metadata"`
	Cells []struct {
		ID        string  `json:"id"`
		RowID     string  `json:"row_id"`
		ColID     string  `json:"col_id"`
		Vertices  struct {
			Xmin float64 `json:"xmin"`
			Xmax float64 `json:"xmax"`
			Ymin float64 `json:"ymin"`
			Ymax float64 `json:"ymax"`
		} `json:"vertices"`
		OcrText   string  `json:"ocr_text"`
		IsHeader  bool    `json:"is_header"`
		Confidence float64 `json:"confidence"`
	} `json:"cells"`
}

func TransformJSON(inputPath, outputPath string) {
	// Read input JSON
	jsonFile, err := os.ReadFile(inputPath)
	if err != nil {
		log.Fatalf("Failed to read getData.json: %v", err)
	}

	var inputData InputTable
	err = json.Unmarshal(jsonFile, &inputData)
	if err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	if len(inputData.Table.Data) == 0 {
		log.Fatalf("No table data found in input JSON.")
	}

	table := inputData.Table.Data[0] // Extract first table entry

	// Create transformed structure
	transformed := TransformedTable{
		Type: "table",
		Metadata: struct {
			ID      string `json:"id"`
			PageNo  int    `json:"page_no"`
			Vertices struct {
				Xmin float64 `json:"xmin"`
				Xmax float64 `json:"xmax"`
				Ymin float64 `json:"ymin"`
				Ymax float64 `json:"ymax"`
			} `json:"vertices"`
		}{
			ID:     table.ID,
			PageNo: table.PageNo,
			Vertices: struct {
				Xmin float64 `json:"xmin"`
				Xmax float64 `json:"xmax"`
				Ymin float64 `json:"ymin"`
				Ymax float64 `json:"ymax"`
			}{Xmin: 0, Xmax: 0, Ymin: 0, Ymax: 0}, // Adjust if required
		},
		Cells: []struct {
			ID        string  `json:"id"`
			RowID     string  `json:"row_id"`
			ColID     string  `json:"col_id"`
			Vertices  struct {
				Xmin float64 `json:"xmin"`
				Xmax float64 `json:"xmax"`
				Ymin float64 `json:"ymin"`
				Ymax float64 `json:"ymax"`
			} `json:"vertices"`
			OcrText   string  `json:"ocr_text"`
			IsHeader  bool    `json:"is_header"`
			Confidence float64 `json:"confidence"`
		}{},
	}

	// Transform cells
	for _, cell := range table.Cells {
		rowID := fmt.Sprintf("row_%d", cell.Row)
		colID := fmt.Sprintf("col_%d", cell.Col)

		transformed.Cells = append(transformed.Cells, struct {
			ID        string  `json:"id"`
			RowID     string  `json:"row_id"`
			ColID     string  `json:"col_id"`
			Vertices  struct {
				Xmin float64 `json:"xmin"`
				Xmax float64 `json:"xmax"`
				Ymin float64 `json:"ymin"`
				Ymax float64 `json:"ymax"`
			} `json:"vertices"`
			OcrText   string  `json:"ocr_text"`
			IsHeader  bool    `json:"is_header"`
			Confidence float64 `json:"confidence"`
		}{
			ID:       cell.ID,
			RowID:    rowID,
			ColID:    colID,
			Vertices: struct {
				Xmin float64 `json:"xmin"`
				Xmax float64 `json:"xmax"`
				Ymin float64 `json:"ymin"`
				Ymax float64 `json:"ymax"`
			}{cell.Xmin, cell.Xmax, cell.Ymin, cell.Ymax},
			OcrText:   cell.Text,
			IsHeader:  true, // Assuming all are headers, adjust if needed
			Confidence: cell.Score,
		})
	}

	// Write transformed JSON
	outputData, err := json.MarshalIndent(transformed, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal transformed JSON: %v", err)
	}

	err = ioutil.WriteFile(outputPath, outputData, 0644)
	if err != nil {
		log.Fatalf("Failed to write transformed JSON: %v", err)
	}

	fmt.Println("✅ Transformation completed! Output saved to", outputPath)
}
