// A simple RESTful API in Go for managing employees	
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

type Employee struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Position string `json:"position"`
}


var (
	employees = make(map[string]Employee)
	mu        sync.Mutex
)

func main() {
	
	mux := http.NewServeMux()

	// Registering endpoints
	mux.HandleFunc("GET /employees", getEmployees)      
	mux.HandleFunc("GET /employees/{id}", getEmployee)  
	mux.HandleFunc("POST /employees", createEmployee)   
	mux.HandleFunc("PUT /employees/{id}", updateEmployee) 
	mux.HandleFunc("DELETE /employees/{id}", deleteEmployee) 

	log.Println("Server starting on :8080...")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}



func getEmployees(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	
	var list []Employee
	for _, emp := range employees {
		list = append(list, emp)
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func getEmployee(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id") 
	
	mu.Lock()
	emp, txst := employees[id]
	mu.Unlock()

	if !txst {
		http.Error(w, "Employee not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(emp)
}

func createEmployee(w http.ResponseWriter, r *http.Request) {
	var emp Employee
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mu.Lock()
	employees[emp.ID] = emp
	mu.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(emp)
}

func updateEmployee(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var updatedEmp Employee
	if err := json.NewDecoder(r.Body).Decode(&updatedEmp); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mu.Lock()
	if _, txst := employees[id]; !txst {
		mu.Unlock()
		http.Error(w, "Employee not found", http.StatusNotFound)
		return
	}
	updatedEmp.ID = id
	employees[id] = updatedEmp
	mu.Unlock()

	json.NewEncoder(w).Encode(updatedEmp)
}

func deleteEmployee(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	mu.Lock()
	delete(employees, id)
	mu.Unlock()

	w.WriteHeader(http.StatusNoContent)
}
