package main

import (
	"go-loan-management-api/internal/config"
	"go-loan-management-api/internal/db"
	"go-loan-management-api/internal/handler"
	"go-loan-management-api/internal/httpx"
	"go-loan-management-api/internal/middleware"
	"go-loan-management-api/internal/repository"
	"go-loan-management-api/internal/service"
	"log"
	"net/http"
)

func main() {
	// Load application configuration from environment variables.
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Open and verify the Postgres connection.
	postgresDB, err := db.NewPostgresConnection(cfg.DBURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer postgresDB.Close()

	log.Println("postgres connection established successfully")

	// Create repositories.
	customerRepository := repository.NewPostgresCustomerRepository(postgresDB)
	loanRepository := repository.NewPostgresLoanRepository(postgresDB)
	repaymentRespository := repository.NewPostgresRepaymentRepository(postgresDB)
	userRepository := repository.NewPostgresUserRepository(postgresDB)

	// Loans and repayments still use in-memory storage for now.
	// memoryStore := store.NewMemoryStore()

	// Services backed by repositories.
	customerService := service.NewCustomerService(customerRepository, loanRepository)
	loanService := service.NewLoanService(customerRepository, loanRepository)
	repaymentsService := service.NewRepaymentService(postgresDB, loanRepository, repaymentRespository)
	authService := service.NewAuthService(
		userRepository,
		cfg.JWTSecret,
		cfg.JWTExpiryHours,
	)

	// Create handlers.
	// 3. Manual Pointer Creation:
	// HealthHandler is a simple struct with no "New" function.
	// We use & to create the pointer manually right here.
	// healthHandler := &handler.HealthHandler{}
	healthHandler := handler.NewHealthHandler(postgresDB)

	// 4. Double Pointer Prevention:
	// customerHandler is already a pointer because NewCustomerHandler returns *CustomerHandler.
	customerHandler := handler.NewCustomerHandler(customerService)
	loansHandler := handler.NewLoanHandler(loanService, repaymentsService)
	repaymentsHandler := handler.NewRepaymentHandler(repaymentsService)
	authHandler := handler.NewAuthHandler(authService)
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret)

	// Create an explicit router instead of using the default global mux.
	mux := http.NewServeMux()

	httpx.RegisterRoutes(mux, healthHandler, customerHandler, loansHandler, repaymentsHandler, authHandler, authMiddleware)

	// Apply middleware in order.
	// RequestID -> Logging -> Recovery -> mux
	var handlerWithMiddleware http.Handler = mux
	handlerWithMiddleware = middleware.Recovery(handlerWithMiddleware)
	handlerWithMiddleware = middleware.Logging(handlerWithMiddleware)
	handlerWithMiddleware = middleware.RequestID(handlerWithMiddleware)

	log.Printf("server running on :%s", cfg.AppPort)
	log.Fatal(http.ListenAndServe(":"+cfg.AppPort, handlerWithMiddleware))
}
