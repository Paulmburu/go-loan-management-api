package httpx

import (
	"go-loan-management-api/internal/handler"
	"go-loan-management-api/internal/middleware"
	"go-loan-management-api/internal/model"
	"go-loan-management-api/internal/response"
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
	authHandler *handler.AuthHandler,
	authMiddleware *middleware.AuthMiddleware,
) {
	adminOrOfficer := middleware.RequireRoles(model.RoleAdmin, model.RoleLoanOfficer)

	mux.HandleFunc("/health", healthHandler.GetHealth)

	// mux.HandleFunc("/customers", HandleByMethod(MethodHandler{
	// 	http.MethodPost: customerHandler.CreateCustomer,
	// 	http.MethodGet:  customerHandler.ListCustomers,
	// }))

	mux.Handle(
		"/customers",
		authMiddleware.RequireAuth(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case http.MethodPost:
					adminOrOfficer(http.HandlerFunc(customerHandler.CreateCustomer)).ServeHTTP(w, r)
				case http.MethodGet:
					adminOrOfficer(http.HandlerFunc(customerHandler.ListCustomers)).ServeHTTP(w, r)
				default:
					response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
				}
			}),
		),
	)

	// mux.HandleFunc("/customers/", func(w http.ResponseWriter, r *http.Request) {
	// 	switch {
	// 	case strings.HasSuffix(r.URL.Path, "/loans"):
	// 		customerHandler.GetCustomerLoans(w, r)
	// 	default:
	// 		customerHandler.GetCustomerByID(w, r)
	// 	}
	// })

	mux.Handle(
		"/customers/",
		authMiddleware.RequireAuth(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch {
				case strings.HasSuffix(r.URL.Path, "/loans"):
					customerHandler.GetCustomerLoans(w, r)
				default:
					customerHandler.GetCustomerByID(w, r)
				}
			}),
		),
	)

	// mux.HandleFunc("/loans", HandleByMethod(MethodHandler{
	// 	http.MethodPost: loanHandler.CreateLoan,
	// 	http.MethodGet:  loanHandler.ListLoans,
	// }))

	mux.Handle(
		"/loans",
		authMiddleware.RequireAuth(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case http.MethodPost:
					adminOrOfficer(http.HandlerFunc(loanHandler.CreateLoan)).ServeHTTP(w, r)
				case http.MethodGet:
					adminOrOfficer(http.HandlerFunc(loanHandler.ListLoans)).ServeHTTP(w, r)
				default:
					response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
				}
			}),
		),
	)

	// mux.HandleFunc("/loans/", func(w http.ResponseWriter, r *http.Request) {
	// 	switch {
	// 	case strings.HasSuffix(r.URL.Path, "/balance"):
	// 		loanHandler.GetLoanBalance(w, r)
	// 	case strings.HasSuffix(r.URL.Path, "/repayments"):
	// 		loanHandler.GetLoanRepayments(w, r)
	// 	default:
	// 		loanHandler.GetLoanByID(w, r)
	// 	}
	// })

	mux.Handle(
		"/loans/",
		authMiddleware.RequireAuth(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch {
				case strings.HasSuffix(r.URL.Path, "/balance"):
					loanHandler.GetLoanBalance(w, r)
				case strings.HasSuffix(r.URL.Path, "/repayments"):
					loanHandler.GetLoanRepayments(w, r)
				default:
					loanHandler.GetLoanByID(w, r)
				}
			}),
		),
	)

	// mux.HandleFunc("/repayments", HandleByMethod(MethodHandler{
	// 	http.MethodPost: repaymentHandler.AddRepayment,
	// }))
	mux.Handle(
		"/repayments",
		authMiddleware.RequireAuth(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case http.MethodPost:
					adminOrOfficer(http.HandlerFunc(repaymentHandler.AddRepayment)).ServeHTTP(w, r)
				default:
					response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
				}
			}),
		),
	)

	mux.HandleFunc("/repayments/", repaymentHandler.GetRepaymentByID)

	mux.HandleFunc("/auth/register", authHandler.Register)
	mux.HandleFunc("/auth/login", authHandler.Login)
	mux.Handle("/auth/me", authMiddleware.RequireAuth(http.HandlerFunc(authHandler.Me)))
}
