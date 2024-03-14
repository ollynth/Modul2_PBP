package main

import (
	"fmt"
	"log"
	cntrl "mod2/controllers"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/users", cntrl.GetAllUsers).Methods("GET")
	router.HandleFunc("/users", cntrl.InsertNewUser).Methods("POST")
	router.HandleFunc("/users", cntrl.UpdateUser).Methods("PUT")
	router.HandleFunc("/users", cntrl.DeleteUser).Methods("DELETE")
	router.HandleFunc("/products", cntrl.GetAllProducts).Methods("GET")
	router.HandleFunc("/products", cntrl.InsertNewProduct).Methods("POST")
	router.HandleFunc("/products", cntrl.UpdateProduct).Methods("PUT")
	router.HandleFunc("/products", cntrl.DeleteProduct).Methods("DELETE")
	router.HandleFunc("/transactions", cntrl.GetAllTransactions).Methods("GET")
	router.HandleFunc("/transactions", cntrl.InsertNewTransaction).Methods("POST")
	router.HandleFunc("/transactions", cntrl.UpdateTransaction).Methods("PUT")
	router.HandleFunc("/products", cntrl.DeleteTransaction).Methods("DELETE")

	http.Handle("/", router)
	fmt.Println("Connected to port 8888")
	log.Println("Connected to port 8888")
	log.Fatal(http.ListenAndServe(":8888", router))
}
