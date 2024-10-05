package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
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

// Handler function for managing individual customers by ID.
func CustomerByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/customers/"):]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		getCustomerByID(w, id)
	case http.MethodPut:
		updateCustomer(w, r, id)
	case http.MethodDelete:
		deleteCustomer(w, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Function to get a customer by ID.
func getCustomerByID(w http.ResponseWriter, id int) {
	for _, customer := range customers {
		if customer.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(customer)
			return
		}
	}
	http.Error(w, "Customer not found", http.StatusNotFound)
}

// Function to update a customer by ID.
func updateCustomer(w http.ResponseWriter, r *http.Request, id int) {
	var updatedCustomer BankCustomer
	if err := json.NewDecoder(r.Body).Decode(&updatedCustomer); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Find the current customer being updated
	var currentCustomer *BankCustomer
	for i := range customers {
		if customers[i].ID == id {
			currentCustomer = &customers[i]
			break
		}
	}

	if currentCustomer == nil {
		http.Error(w, "Customer not found", http.StatusNotFound)
		return
	}

	// Check for unique email only if it's different from the current email
	if updatedCustomer.Email != currentCustomer.Email && !isEmailUnique(updatedCustomer.Email) {
		http.Error(w, "Email already exists", http.StatusConflict)
		return
	}

	for i, customer := range customers {
		if customer.ID == id {
			// Update the fields of the existing customer
			customers[i].Name = updatedCustomer.Name
			customers[i].Email = updatedCustomer.Email
			customers[i].Balance = updatedCustomer.Balance
			customers[i].AccountType = updatedCustomer.AccountType

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(customers[i])
			return
		}
	}
	http.Error(w, "Customer not found", http.StatusNotFound)
}

// Function to delete a customer by ID.
func deleteCustomer(w http.ResponseWriter, id int) {
	for i, customer := range customers {
		if customer.ID == id {
			customers = append(customers[:i], customers[i+1:]...) // Remove the customer

			w.Header().Set("Content-Type", "application/json")
			response := map[string]string{"message": "Customer deleted successfully"}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
			return
		}
	}
	http.Error(w, "Customer not found", http.StatusNotFound)
}

func main() {

	http.HandleFunc("/customers", CustomersHandler)     // For managing multiple customers
	http.HandleFunc("/customers/", CustomerByIDHandler) // For managing individual customers by ID
	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
