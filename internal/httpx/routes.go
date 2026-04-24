package httpx

import (
	"go-loan-management-api/internal/handler"
	"net/http"
	"strings"
)

// RegisterRoutes registers all HTTP routes on the mux.
func RegisterRoutes(
	mux *http.ServeMux,
	healthHandler *handler.HealthHandler,
	customerHandler *handler.CustomerHandler,
	loanHandler *handler.LoanHandler,
	repaymentHandler *handler.RepaymentHandler,
) {
	mux.HandleFunc("/health", healthHandler.GetHealth)

	mux.HandleFunc("/customers", HandleByMethod(MethodHandler{
		http.MethodPost: customerHandler.CreateCustomer,
		http.MethodGet:  customerHandler.ListCustomers,
	}))

	mux.HandleFunc("/customers/", func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/loans"):
			customerHandler.GetCustomerLoans(w, r)
		default:
			customerHandler.GetCustomerByID(w, r)
		}
	})

	mux.HandleFunc("/loans", HandleByMethod(MethodHandler{
		http.MethodPost: loanHandler.CreateLoan,
		http.MethodGet:  loanHandler.ListLoans,
	}))

	mux.HandleFunc("/loans/", func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/balance"):
			loanHandler.GetLoanBalance(w, r)
		case strings.HasSuffix(r.URL.Path, "/repayments"):
			loanHandler.GetLoanRepayments(w, r)
		default:
			loanHandler.GetLoanByID(w, r)
		}
	})

	mux.HandleFunc("/repayments", HandleByMethod(MethodHandler{
		http.MethodPost: repaymentHandler.AddRepayment,
	}))

	mux.HandleFunc("/repayments/", repaymentHandler.GetRepaymentByID)
}
