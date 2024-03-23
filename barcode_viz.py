import os
import json
import pandas as pd

# Define the directory where the JSON files are stored
output_dir = "checkout-items"

# Initialize a list to store the product data
data = []

# Loop over all the JSON files in the directory
for filename in os.listdir(output_dir):
    if filename.endswith("-description.json"):
        # Open the JSON file and load the data
        with open(os.path.join(output_dir, filename), "r") as json_file:
            product_info = json.load(json_file)

        # Loop over the items in the product info
        for item in product_info.get('items', []):
            # Append the title, brand, model, and category to the data list
            data.append({
                'Title': item.get('title'),
                'Brand': item.get('brand'),
                'Model': item.get('model'),
                'Category': item.get('category')
            })

# Create a pandas DataFrame from the data
df = pd.DataFrame(data)

# Convert the DataFrame to an HTML table
table_html = df.to_html(index=False)

# Create the HTML page
html = f"""
<html>
<head>
    <style>
        body {{
            font-family: Arial, sans-serif;
            background-color: #333;
            color: #fff;
            padding: 20px;
            display: flex;
            flex-direction: column;
            align-items: center;
        }}
        .container {{
            width: 70%;
        }}
        table {{
            width: 100%;
            border-collapse: collapse;
        }}
        th, td {{
            border: 1px solid #ddd;
            padding: 8px;
        }}
        th {{
            background-color: #4CAF50;
            color: white;
        }}
    </style>
</head>
<body>
    <div class="container">
        <h1>🛒 Full-Stack Checkout Analysis 🛒</h1>
        <p>This project provides a full-stack overview of a checkout process at a grocery store. It involves a barcode scanner that is connected to your local machine. When you scan a product's barcode, the scanner sends the UPC (Universal Product Code) to your computer. A script then takes this UPC, makes a request to the UPCitemdb API to get the product information, and saves this information to a JSON file in the "checkout-items" directory.</p>
        <p>This HTML page reads those JSON files, extracts the product information, and displays it in the following data table:</p>
        {table_html}
    </div>
</body>
</html>
"""

# Write the HTML to a file
with open("inventory.html", "w", encoding='utf-8') as file:
    file.write(html)