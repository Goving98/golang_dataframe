import json

def transform_json(input_json):
    table_data = input_json.get("table", {}).get("data", [])[0]
    
    transformed = {
        "type": "table",
        "metadata": {
            "id": "table_1",
            "page_no": table_data.get("page_no", 0),
            "vertices": {
                "xmin": table_data.get("xmin", 0),
                "xmax": table_data.get("xmax", 0),
                "ymin": table_data.get("ymin", 0),
                "ymax": table_data.get("ymax", 0)
            }
        },
        "cells": [],
        "layout": {"row_order": [], "column_order": []},
        "headers": {}
    }
    
    column_map = {}
    row_ids = set()
    
    for cell in table_data.get("cells", []):
        cell_id = cell["id"]
        row_id = str(cell["row"])
        col_id = str(cell["col"])
        row_ids.add(row_id)
        
        if cell.get("label") and row_id == "1":
            column_map[col_id] = cell["label"]
        
        transformed["cells"].append({
            "id": cell_id,
            "row_id": row_id,
            "col_id": col_id,
            "vertices": {
                "xmin": cell["xmin"],
                "xmax": cell["xmax"],
                "ymin": cell["ymin"],
                "ymax": cell["ymax"]
            },
            "ocr_text": cell.get("text", ""),
            "is_header": row_id == "1",
            "confidence": cell.get("score", 0.0)
        })
    
    transformed["layout"]["row_order"] = sorted(row_ids, key=int)
    transformed["layout"]["column_order"] = sorted(column_map.keys(), key=int)
    transformed["headers"] = {f"col_{col}": label for col, label in column_map.items()}
    
    return transformed

# Example usage
if __name__ == "__main__":
    with open("getData.json", "r") as file:
        input_data = json.load(file)
    
    output_data = transform_json(input_data)
    
    with open("finalData.json", "w") as file:
        json.dump(output_data, file, indent=2)
    
    print("Transformation complete. Output saved to output.json")