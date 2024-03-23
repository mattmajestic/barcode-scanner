// main.go

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Product struct {
	Title    string `json:"title"`
	Brand    string `json:"brand"`
	Model    string `json:"model"`
	Category string `json:"category"`
	Price    string `json:"price"`
}

type Response struct {
	Items []Product `json:"items"`
}

func main() {
	fmt.Println("\033[1;36mListening for barcode scans. Press CTRL+C to exit.\033[0m")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		barcodeData := scanner.Text()

		// Print the received barcode
		fmt.Printf("\033[1;32mReceived barcode: %s\033[0m\n", barcodeData)

		// Make a GET request to the UPCitemdb API to get the product information
		resp, err := http.Get(fmt.Sprintf("https://api.upcitemdb.com/prod/trial/lookup?upc=%s", barcodeData))
		if err != nil {
			log.Fatalf("Error making GET request: %v", err)
		}

		// Parse the JSON response
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Error reading response body: %v", err)
		}
		var productInfo Response
		err = json.Unmarshal(body, &productInfo)
		if err != nil {
			log.Fatalf("Error unmarshalling JSON: %v", err)
		}

		// Print the product information in a more structured way
		for _, item := range productInfo.Items {
			fmt.Printf("\033[1;34mProduct Title: %s\033[0m\n", item.Title)
			fmt.Printf("\033[1;34mProduct Brand: %s\033[0m\n", item.Brand)
			fmt.Printf("\033[1;34mProduct Model: %s\033[0m\n", item.Model)
			fmt.Printf("\033[1;34mProduct Category: %s\033[0m\n", item.Category)
			fmt.Printf("\033[1;34mProduct Price: %s\033[0m\n", item.Price)
		}
	}
}
