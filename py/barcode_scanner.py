import os
import requests
import json
from colorama import init, Fore, Style

# Initialize Colorama
init(autoreset=True)

# Define the directory and file path for output
output_dir = "checkout-items"
output_file = "scanned_barcodes.txt"
full_path = os.path.join(output_dir, output_file)

# Ensure the output directory exists
os.makedirs(output_dir, exist_ok=True)

print(Fore.YELLOW + Style.BRIGHT + "Listening for barcode scans. Press CTRL+C to exit.")

try:
    with open(full_path, "a") as file:
        while True:
            # Simulate listening for barcode scanner input (from standard input)
            barcode_data = input(Fore.CYAN + "Scan a barcode: ")
            
            # Print the received barcode in green
            print(Fore.GREEN + Style.BRIGHT + f"Received barcode: {barcode_data}")
            
            # Write the scanned barcode data to the file
            file.write(barcode_data + "\n")

            # Make a GET request to the UPCitemdb API to get the product information
            response = requests.get(f"https://api.upcitemdb.com/prod/trial/lookup?upc={barcode_data}")

            # Check if the request was successful
            if response.status_code == 200:
                # Parse the JSON response
                product_info = response.json()

                # Print the product information in a more structured way
                for item in product_info.get('items', []):
                    print(Fore.GREEN + Style.BRIGHT + f"Product Title: {item.get('title')}")
                    print(Fore.GREEN + Style.BRIGHT + f"Product Brand: {item.get('brand')}")
                    print(Fore.GREEN + Style.BRIGHT + f"Product Model: {item.get('model')}")
                    print(Fore.GREEN + Style.BRIGHT + f"Product Category: {item.get('category')}")

                # Save the UPC and associated JSON data to a file
                with open(f"{output_dir}/{barcode_data}-description.json", "w") as json_file:
                    json.dump(product_info, json_file)
            else:
                print(Fore.RED + Style.BRIGHT + f"Failed to get product info for UPC {barcode_data}")

except KeyboardInterrupt:
    print(Fore.RED + Style.BRIGHT + "Exiting barcode listener.")