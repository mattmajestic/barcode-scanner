package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Product struct {
	Title string `json:"title"`
	Brand string `json:"brand"`
	Model string `json:"model"`
}

type Response struct {
	Items []Product `json:"items"`
}

var db *sql.DB

func init() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connStr := os.Getenv("NEON_URL")
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	http.HandleFunc("/", productHandler)
	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	// Setup Scanner to read barcode data from stdin
	scanner := bufio.NewScanner(os.Stdin)

	// Create a slice to store the products
	var products []Product

	for {
		// Print "Ready to scan" in blue
		fmt.Println("\033[1;34mReady to scan\033[0m")

		if !scanner.Scan() {
			break
		}

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

		// Print the product information and add it to the products slice
		for _, item := range productInfo.Items {
			fmt.Printf("\033[1;34mProduct Title: %s\033[0m\n", item.Title)
			fmt.Printf("\033[1;34mProduct Brand: %s\033[0m\n", item.Brand)
			fmt.Printf("\033[1;34mProduct Model: %s\033[0m\n", item.Model)

			// Add the product to the slice
			products = append(products, Product{Title: item.Title, Brand: item.Brand, Model: item.Model})

			// Insert the product information into the database
			_, err := db.Exec(`INSERT INTO product (title, brand, model) VALUES ($1, $2, $3)`, item.Title, item.Brand, item.Model)
			if err != nil {
				log.Fatalf("Error inserting data into database: %v", err)
			}
		}
	}

	// Print the products
	for _, product := range products {
		fmt.Printf("Title: %s, Brand: %s, Model: %s\n", product.Title, product.Brand, product.Model)
	}
}

func productHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT title, brand, model FROM product")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.Title, &product.Brand, &product.Model); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		products = append(products, product)
	}

	tmpl, err := template.ParseFiles("templates/products.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, products)
}
