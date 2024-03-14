package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	m "mod2/models"
	"net/http"
	"strings"
)
// success response
func sendSuccessResponseUser(w http.ResponseWriter, message string) {
	var response m.UserResponse
	response.Status = 200
	response.Message = message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// get all user's data
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	query := "SELECT * FROM users"

	name := r.URL.Query()["name"]
	age := r.URL.Query()["age"]
	if name != nil {
		fmt.Println(name[0])
		query += " WHERE name='" + name[0] + "'"
	}
	if age != nil {
		if name[0] != "" {
			query += " AND"
		} else {
			query += " WHERE"
		}
		query += " age='" + age[0] + "'"
	}

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		return
	}
	var user m.Users
	var users []m.Users
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Name, &user.Age, &user.Address); err != nil {
			log.Println(err)
		} else {
			users = append(users, user)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	var response m.UsersResponse
	response.Status = 200
	response.Message = "success"
	response.Data = users
	json.NewEncoder(w).Encode(response)
}

// insert new user
func InsertNewUser(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		log.Println("Error parsing form data:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if r.Form.Get("name") == "" || r.Form.Get("age") == "" || r.Form.Get("address") == "" {
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

	// buat query
	query := "INSERT INTO users (name, age, address) VALUES (?, ?, ?)"
	stmt, err := tx.Prepare(query)
	if err != nil {
		log.Println("Error preparing SQL statement:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(r.Form.Get("name"), r.Form.Get("age"), r.Form.Get("address"))
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

	fmt.Fprintf(w, "User inserted successfully with ID: %d", lastInsertID)
}

// update user's data
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	// get userID
	userID := r.URL.Query().Get("id")

	// Validate user ID
	if userID == "" {
		log.Println("Error: Missing user ID")
		http.Error(w, "Bad Request: Missing user ID", http.StatusBadRequest)
		return
	}

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		log.Println("Error parsing form data:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Get updated values from form data
	name := r.Form.Get("name")
	age := r.Form.Get("age")
	address := r.Form.Get("address")

	// Validate updated values
	if name == "" && age == "" && address == "" {
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

	// Construct the update query based on the provided values
	query := "UPDATE users SET"
	var updates []string

	if name != "" {
		updates = append(updates, "name = ?")
	}

	if age != "" {
		updates = append(updates, "age = ?")
	}

	if address != "" {
		updates = append(updates, "address = ?")
	}

	//gabungin yang mau diupdate
	updateString := strings.Join(updates, ", ")

	// Add the WHERE clause for the specific user ID
	query += " " + updateString + " WHERE id = ?"

	stmt, err := tx.Prepare(query)
	if err != nil {
		log.Println("Error preparing SQL statement:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	// Execute the update query
	_, err = stmt.Exec(name, age, address, userID)
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

	fmt.Fprintf(w, "User with ID %s updated successfully", userID)
}

// delete user
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	userID := r.URL.Query().Get("id")

	// cek user id yang mau dihapus
	if userID == "" {
		log.Println("Error: User ID tidak ada")
		http.Error(w, "Bad Request: User ID tidak ada", http.StatusBadRequest)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Println("Error :", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// delete query
	query := "DELETE FROM users WHERE id = ?"

	stmt, err := tx.Prepare(query)
	if err != nil {
		log.Println("AQL statement Error: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	// Execute the delete query
	_, err = stmt.Exec(userID)
	if err != nil {
		log.Println("Error executing SQL statement:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "User with ID %s deleted successfully", userID)
}
