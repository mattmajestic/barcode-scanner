// main.go

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Product struct {
	Title    string `json:"title"`
	Brand    string `json:"brand"`
	Model    string `json:"model"`
	Category string `json:"category"`
}

type Response struct {
	Items []Product `json:"items"`
}

func main() {
	outputDir := "checkout-items"
	outputFile := "scanned_barcodes.txt"
	fullPath := outputDir + "/" + outputFile

	// Ensure the output directory exists
	os.MkdirAll(outputDir, os.ModePerm)

	fmt.Println("Listening for barcode scans. Press CTRL+C to exit.")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		barcodeData := scanner.Text()

		// Print the received barcode
		fmt.Printf("Received barcode: %s\n", barcodeData)

		// Write the scanned barcode data to the file
		f, _ := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		f.WriteString(barcodeData + "\n")
		f.Close()

		// Make a GET request to the UPCitemdb API to get the product information
		resp, _ := http.Get(fmt.Sprintf("https://api.upcitemdb.com/prod/trial/lookup?upc=%s", barcodeData))

		// Check if the request was successful
		if resp.StatusCode == 200 {
			// Parse the JSON response
			body, _ := ioutil.ReadAll(resp.Body)
			var productInfo Response
			json.Unmarshal(body, &productInfo)

			// Print the product information in a more structured way
			for _, item := range productInfo.Items {
				fmt.Printf("Product Title: %s\n", item.Title)
				fmt.Printf("Product Brand: %s\n", item.Brand)
				fmt.Printf("Product Model: %s\n", item.Model)
				fmt.Printf("Product Category: %s\n", item.Category)
			}

			// Save the UPC and associated JSON data to a file
			jsonFile, _ := json.MarshalIndent(productInfo, "", " ")
			ioutil.WriteFile(fmt.Sprintf("%s/%s-description.json", outputDir, barcodeData), jsonFile, 0644)
		} else {
			fmt.Printf("Failed to get product info for UPC %s\n", barcodeData)
		}
	}
}