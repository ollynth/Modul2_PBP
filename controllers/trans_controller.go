package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	m "mod2/models"
	"net/http"
	"strings"
)

func GetAllTransactions(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	query := "SELECT * FROM transactions"

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		return
	}
	var transaction m.Transactions
	var transactions []m.Transactions
	for rows.Next() {
		if err := rows.Scan(&transaction.ID, &transaction.UserID, &transaction.ProductID, &transaction.Quantity); err != nil {
			log.Println(err)
		} else {
			transactions = append(transactions, transaction)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	var response m.TransactionsResponse
	response.Status = 200
	response.Message = "success"
	response.Data = transactions
	json.NewEncoder(w).Encode(response)
}

// buat insert trans
func InsertNewTransaction(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		log.Println("Error parsing form data:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if r.Form.Get("userid") == "" || r.Form.Get("productid") == "" || r.Form.Get("quantity") == "" {
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

	query := "INSERT INTO products (userid, productid, quantity) VALUES (?, ?, ?)"
	stmt, err := tx.Prepare(query)
	if err != nil {
		log.Println("Error preparing SQL statement:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(r.Form.Get("userid"), r.Form.Get("productid"), r.Form.Get("quantity"))
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

	fmt.Fprintf(w, "Transaction inserted successfully with ID: %d", lastInsertID)
}

func UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	transactionID := r.URL.Query().Get("id")

	if transactionID == "" {
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

	userid := r.Form.Get("userid")
	productid := r.Form.Get("productid")
	quantity := r.Form.Get("quantity")

	if userid == "" && productid == "" && quantity == "" {
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

	query := "UPDATE transactions SET"
	var updates []string

	if userid != "" {
		updates = append(updates, "userid = ?")
	}

	if productid != "" {
		updates = append(updates, "productid = ?")
	}

	if quantity != "" {
		updates = append(updates, "quantity = ?")
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

	_, err = stmt.Exec(userid, productid, quantity)
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

	fmt.Fprintf(w, "Transaction with ID %s updated successfully", transactionID)
}

func DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	transactionID := r.URL.Query().Get("id")

	if transactionID == "" {
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

	query := "DELETE FROM transactions WHERE id = ?"

	stmt, err := tx.Prepare(query)
	if err != nil {
		log.Println("Error preparing SQL statement:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(transactionID)
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

	fmt.Fprintf(w, "Transaction with ID %s deleted successfully", transactionID)
}
