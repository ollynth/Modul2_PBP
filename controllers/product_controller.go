package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	m "mod2/models"
	"net/http"
	"strings"
)

func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	query := "SELECT * FROM products"

	name := r.URL.Query()["name"]
	price := r.URL.Query()["price"]
	if name != nil {
		fmt.Println(name[0])
		query += " WHERE name='" + name[0] + "'"
	}
	if price != nil {
		if name[0] != "" {
			query += " AND"
		} else {
			query += " WHERE"
		}
		query += " price='" + price[0] + "'"
	}

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		return
	}
	var product m.Products
	var products []m.Products
	for rows.Next() {
		if err := rows.Scan(&product.ID, &product.Name, &product.Price); err != nil {
			log.Println(err)
		} else {
			products = append(products, product)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	var response m.ProductsResponse
	response.Status = 200
	response.Message = "success"
	response.Data = products
	json.NewEncoder(w).Encode(response)
}

// buat insert produk baru
func InsertNewProduct(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		log.Println("Error parsing form data:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if r.Form.Get("name") == "" || r.Form.Get("price") == "" {
		log.Println("Error: Incomplete data provided")
		http.Error(w, "Bad Request: Incomplete data", http.StatusBadRequest)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Println("Error starting transaction:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	query := "INSERT INTO products (name, price) VALUES (?, ?)"
	stmt, err := tx.Prepare(query)
	if err != nil {
		log.Println("Error preparing SQL statement:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(r.Form.Get("name"), r.Form.Get("price"))
	if err != nil {
		log.Println("Error executing SQL statement:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	lastInsertID, _ := result.LastInsertId()

	err = tx.Commit()
	if err != nil {
		log.Println("Error committing transaction:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Product inserted successfully with ID: %d", lastInsertID)
}

// buat update data produk
func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	productID := r.URL.Query().Get("id")

	if productID == "" {
		log.Println("Error: Missing product ID")
		http.Error(w, "Bad Request: Missing product ID", http.StatusBadRequest)
		return
	}

	err := r.ParseForm()
	if err != nil {
		log.Println("Error parsing form data:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	name := r.Form.Get("name")
	price := r.Form.Get("price")

	if name == "" && price == "" {
		log.Println("Error: No data provided for update")
		http.Error(w, "Bad Request: No data provided for update", http.StatusBadRequest)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Println("Error starting transaction:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	query := "UPDATE products SET"
	var updates []string

	if name != "" {
		updates = append(updates, "name = ?")
	}

	if price != "" {
		updates = append(updates, "price = ?")
	}

	updateString := strings.Join(updates, ", ")

	query += " " + updateString + " WHERE id = ?"

	stmt, err := tx.Prepare(query)
	if err != nil {
		log.Println("Error preparing SQL statement:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(name, price, productID)
	if err != nil {
		log.Println("Error executing SQL statement:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Println("Error committing transaction:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Product with ID %s updated successfully", productID)
}

// buat hapus produk
func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	productID := r.URL.Query().Get("id")

	if productID == "" {
		log.Println("Error: Missing product ID")
		http.Error(w, "Bad Request: Missing product ID", http.StatusBadRequest)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Println("Error starting transaction:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	query := "DELETE FROM products WHERE id = ?"

	stmt, err := tx.Prepare(query)
	if err != nil {
		log.Println("Error preparing SQL statement:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(productID)
	if err != nil {
		log.Println("Error executing SQL statement:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Println("Error committing transaction:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Product with ID %s deleted successfully", productID)
}
