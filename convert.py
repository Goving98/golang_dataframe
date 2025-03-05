import json

# Load check.json (unwanted format)
check_json_path = "original.json"
with open(check_json_path, "r", encoding="utf-8") as file:
    check_data = json.load(file)

# Extract relevant table data from check.json
table_data = check_data["table"]["data"][0]

# Prepare the transformed JSON in the desired format
transformed_json = {
    "type": "table",
    "metadata": {
        "id": "table_1",
        "page_no": table_data["page_no"],
        "vertices": {
            "xmin": table_data["xmin"],
            "xmax": table_data["xmax"],
            "ymin": table_data["ymin"],
            "ymax": table_data["ymax"]
        }
    },
    "cells": [],
    "layout": {
        "row_order": [],
        "column_order": []
    },
    "headers": {}
}


row_order = set()
column_order = set()

for cell in table_data["cells"]:
    row_id = f"row_{cell['row']}"
    col_id = f"col_{cell['col']}"
    
    transformed_cell = {
        "id": f"cell_{cell['row']}_{cell['col']}",
        "row_id": row_id,
        "col_id": col_id,
        "vertices": {
            "xmin": cell["xmin"],
            "xmax": cell["xmax"],
            "ymin": cell["ymin"],
            "ymax": cell["ymax"]
        },
        "ocr_text": str(cell["text"]),  # Convert to string
        "is_header": cell["label"] == "header",  # Define header logic
        "confidence": cell["score"]
    }

    transformed_json["cells"].append(transformed_cell)
    row_order.add(row_id)
    column_order.add(col_id)

    # Map headers based on the first row of data
    if cell["row"] == 1:
        transformed_json["headers"][col_id] = cell["label"]

# Sort and update row/column order
transformed_json["layout"]["row_order"] = sorted(row_order)
transformed_json["layout"]["column_order"] = sorted(column_order)

# Save the transformed JSON to a file
output_json_path = "getData.json"
with open(output_json_path, "w", encoding="utf-8") as file:
    json.dump(transformed_json, file, indent=4)

# Return the path of the transformed JSON
output_json_path
