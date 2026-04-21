package main

import (
	"go-loan-management-api/internal/handler"
	"go-loan-management-api/internal/service"
	"go-loan-management-api/internal/store"
	"log"
	"net/http"
	"strings"
)

func main() {

	// Create the in-memory store shared by all services.
	//
	// 1. Return-based Pointer:
	// NewMemoryStore() already returns a *MemoryStore (pointer).
	// The & is "hidden" inside the function's return statement.
	memoryStore := store.NewMemoryStore()

	// Create services.
	//
	// 2. Dependency Injection:
	// We pass the memoryStore pointer so the service shares the same data address.
	customerService := service.NewCustomerService(memoryStore)
	loanService := service.NewLoanService(memoryStore)
	repaymentsService := service.NewRepaymentService(memoryStore)

	// Create handlers.
	// 3. Manual Pointer Creation:
	// HealthHandler is a simple struct with no "New" function.
	// We use & to create the pointer manually right here.
	healthHandler := &handler.HealthHandler{}

	// 4. Double Pointer Prevention:
	// customerHandler is already a pointer because NewCustomerHandler returns *CustomerHandler.
	customerHandler := handler.NewCustomerHandler(customerService)
	loansHandler := handler.NewLoanHandler(loanService)
	repaymentsHandler := handler.NewRepaymentHandler(repaymentsService)

	// Register routes.
	http.HandleFunc("/health", healthHandler.GetHealth)
	http.HandleFunc("/createCustomer", customerHandler.CreateCustomer)

	// This is an Anonymous Function used as a Method Dispatcher.
	// It acts as a wrapper to route requests based on the HTTP Method (GET vs POST).
	http.HandleFunc("/loans", func(w http.ResponseWriter, r *http.Request) {
		// This function is also a Closure: it captures 'loansHandler' from the outer scope.
		switch r.Method {
		case http.MethodPost:
			loansHandler.CreateLoan(w, r)
		case http.MethodGet:
			loansHandler.ListLoans(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Handle /loans/{id}/balance.
	http.HandleFunc("/loans/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/balance") {
			loansHandler.GetLoanBalance(w, r)
			return
		}
		http.Error(w, "not found", http.StatusNotFound)
	})

	http.HandleFunc("/repayments", repaymentsHandler.AddRepayment)

	log.Println("server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
