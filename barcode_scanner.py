import os
from colorama import init, Fore, Style

# Initialize Colorama
init(autoreset=True)

# Define the directory and file path for output
output_dir = "outputs"
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
except KeyboardInterrupt:
    print(Fore.RED + Style.BRIGHT + "Exiting barcode listener.")
