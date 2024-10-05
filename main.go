package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// BankCustomer represents a customer in the banking system.
type BankCustomer struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Email       string  `json:"email"`
	Balance     float64 `json:"balance"`
	AccountType string  `json:"account_type"`
}

// To store bank customers.
var customers []BankCustomer

// Handler function for managing multiple customers.
func CustomersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		createCustomer(w, r)
	case http.MethodGet:
		getAllCustomers(w)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Function to create a new customer.
func createCustomer(w http.ResponseWriter, r *http.Request) {
	var newCustomer BankCustomer
	if err := json.NewDecoder(r.Body).Decode(&newCustomer); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Check for unique email
	if !isEmailUnique(newCustomer.Email) {
		http.Error(w, "Email already exists", http.StatusConflict)
		return
	}

	newCustomer.ID = len(customers) + 1
	customers = append(customers, newCustomer)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newCustomer)
}

// Function to check if the email is unique.
func isEmailUnique(email string) bool {
	for _, customer := range customers {
		if customer.Email == email {
			return false // Email already exists
		}
	}
	return true // Email is unique
}

// Function to get all customers.
func getAllCustomers(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customers)
}

func main() {

	http.HandleFunc("/customers", CustomersHandler) // For managing multiple customers
	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
