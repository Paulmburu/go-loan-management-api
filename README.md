# Go Loan Management API

A Go backend project for learning backend fundamentals through a loan management domain.

The API supports:

- customers
- loans
- repayments
- loan balances
- request logging
- request IDs
- panic recovery
- Postgres persistence

This project is being built in phases:

- Phase 1: beginner / foundation
- Phase 2: intermediate / persistence + middleware
- Phase 3: advanced / auth + testing + stronger architecture

---

## Features

Current features include:

### Customers

- create customer
- list customers
- get customer by ID
- get loans for a customer

### Loans

- create loan
- list loans
- get loan by ID
- get loan balance
- get repayments for a loan

### Repayments

- add repayment
- get repayment by ID

### Infrastructure

- Postgres connection
- SQL migrations
- request ID middleware
- request logging middleware
- recovery middleware
- DB-aware health check
- better error mapping
- transaction-safe repayment flow
- consistent JSON API responses

---

## Tech stack

- Go
- net/http
- PostgreSQL
- database/sql
- lib/pq

---

## Project structure

```text
go-loan-management-api/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── apperror/
│   ├── config/
│   ├── db/
│   ├── handler/
│   ├── middleware/
│   ├── model/
│   ├── repository/
│   ├── response/
│   └── service/
├── migrations/
├── .env
├── go.mod
└── README.md
```

---

## Prerequisites

Make sure you have installed:

- Go
- PostgreSQL
- psql (or another Postgres client)

---

## Environment variables

Set the following environment variables before running the app:

```env
APP_PORT=8080
DATABASE_URL=postgres://postgres:postgres@localhost:5432/go_loan_management_api?sslmode=disable
```

You can export them manually:

```bash
export APP_PORT=8080
export DATABASE_URL='postgres://postgres:postgres@localhost:5432/go_loan_management_api?sslmode=disable'
```

---

## Database setup

Create the database:

```bash
createdb go_loan_management_api
```

If needed, specify your Postgres user:

```bash
createdb -U postgres go_loan_management_api
```

---

## Running migrations

Run migrations in order:

```bash
psql -d go_loan_management_api -f migrations/001_create_customers.sql
psql -d go_loan_management_api -f migrations/002_create_loans.sql
psql -d go_loan_management_api -f migrations/003_create_repayments.sql
psql -d go_loan_management_api -f migrations/004_add_customer_email_unique_constraint.sql
psql -d go_loan_management_api -f migrations/005_add_loan_check_constraints.sql
psql -d go_loan_management_api -f migrations/006_add_repayment_check_constraints.sql
```

If needed:

```bash
psql -U postgres -d go_loan_management_api -f migrations/001_create_customers.sql
```

Repeat for the remaining files.

---

## Running the app

Install dependencies:

```bash
go mod tidy
```

Start the API:

```bash
go run ./cmd/api
```

Expected startup logs:

```text
postgres connection established successfully
server running on :8080
```

---

## Available endpoints

### Health

- `GET /health`

### Customers

- `POST /customers`
- `GET /customers`
- `GET /customers/{id}`
- `GET /customers/{id}/loans`

### Loans

- `POST /loans`
- `GET /loans`
- `GET /loans/{id}`
- `GET /loans/{id}/balance`
- `GET /loans/{id}/repayments`

### Repayments

- `POST /repayments`
- `GET /repayments/{id}`

---

## Sample requests

### Health check

```bash
curl http://localhost:8080/health
```

### Create customer

```bash
curl -X POST http://localhost:8080/customers \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "Paul Mburu",
    "email": "paul@example.com",
    "phone": "0712345678"
  }'
```

### List customers

```bash
curl http://localhost:8080/customers
```

### Get customer by ID

```bash
curl http://localhost:8080/customers/1
```

### Create loan

```bash
curl -X POST http://localhost:8080/loans \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": 1,
    "principal_amount": 10000,
    "interest_rate": 12
  }'
```

### List loans

```bash
curl http://localhost:8080/loans
```

### Get loan by ID

```bash
curl http://localhost:8080/loans/1
```

### Get loan balance

```bash
curl http://localhost:8080/loans/1/balance
```

### Add repayment

```bash
curl -X POST http://localhost:8080/repayments \
  -H "Content-Type: application/json" \
  -d '{
    "loan_id": 1,
    "amount": 2000
  }'
```

### Get repayment by ID

```bash
curl http://localhost:8080/repayments/1
```

### Get repayments for a loan

```bash
curl http://localhost:8080/loans/1/repayments
```

### Get loans for a customer

```bash
curl http://localhost:8080/customers/1/loans
```

---

## Current architecture

The app currently follows this flow:

```text
request -> middleware -> handler -> service -> repository -> Postgres -> response
```

Responsibilities:

- handler: HTTP request/response logic
- service: business rules
- repository: database access
- middleware: cross-cutting concerns
- response: consistent JSON API formatting

---

## Middleware

Current middleware includes:

- request ID middleware
- logging middleware
- recovery middleware

The middleware setup currently provides:

- `X-Request-ID` on responses
- request ID in request context
- request logging with:
  - request ID
  - method
  - path
  - status
  - duration
- panic recovery with:
  - request ID
  - method
  - path
  - panic value

---

## Error handling

The project uses application-level error mapping for clearer API responses.

Examples:

- not found -> `404`
- conflict -> `409`
- validation issue -> `400`
- internal error -> `500`

This keeps raw database errors from leaking directly to API clients.

---

## Applied backend and system design concepts

This project already applies several important backend and system design concepts, even before Phase 3.

### ACID transactions

ACID principles are currently applied in the **repayment flow**.

When a repayment is processed, the system:

1. inserts the repayment record
2. updates the related loan outstanding balance and status

These operations are executed inside a **single database transaction** so they either both succeed or both fail.

Why this matters:

- prevents a repayment row from being saved without the loan balance being updated
- prevents partial writes during failures
- protects financial consistency

This mainly demonstrates:

- **Atomicity**: the repayment write and loan update happen as one unit
- **Consistency**: the database moves from one valid state to another
- **Isolation** and **Durability** are provided by PostgreSQL transaction behavior

Example snippet from the repayment flow:

```go
tx, err := s.db.BeginTx(ctx, nil)
if err != nil {
	return model.Repayment{}, fmt.Errorf("failed to begin transaction: %w", err)
}
defer func() {
	_ = tx.Rollback()
}()

err = tx.QueryRowContext(ctx, createRepaymentQuery, loanID, amount).Scan(
	&savedRepayment.ID,
	&savedRepayment.LoanID,
	&savedRepayment.Amount,
	&savedRepayment.CreatedAt,
)
if err != nil {
	return model.Repayment{}, fmt.Errorf("failed to create repayment: %w", err)
}

_, err = tx.ExecContext(ctx, updateLoanQuery, loan.OutstandingAmount, loan.Status, loan.ID)
if err != nil {
	return model.Repayment{}, fmt.Errorf("failed to update loan after repayment: %w", err)
}

if err := tx.Commit(); err != nil {
	return model.Repayment{}, fmt.Errorf("failed to commit transaction: %w", err)
}
```

### Data integrity constraints

The database applies integrity rules using:

- unique constraints
- check constraints
- foreign keys

Examples in this project:

- customer email must be unique
- repayment amount must be positive
- principal amount must be positive
- loan status must be one of the allowed values
- loans must reference valid customers
- repayments must reference valid loans

This is a system design concept because correctness is enforced not only in Go code, but also at the persistence layer.

Example SQL constraints:

```sql
ALTER TABLE customers
ADD CONSTRAINT customers_email_unique UNIQUE (email);

ALTER TABLE loans
ADD CONSTRAINT loans_principal_amount_positive CHECK (principal_amount > 0),
ADD CONSTRAINT loans_interest_rate_non_negative CHECK (interest_rate >= 0),
ADD CONSTRAINT loans_status_valid CHECK (status IN ('active', 'paid'));

ALTER TABLE repayments
ADD CONSTRAINT repayments_amount_positive CHECK (amount > 0);
```

### Layered architecture

The project uses a layered backend structure:

```text
request -> middleware -> handler -> service -> repository -> database -> response
```

Where:

- **handlers** deal with HTTP concerns
- **services** hold business rules
- **repositories** isolate persistence logic
- **middleware** handles cross-cutting concerns

This is important because it keeps responsibilities separated and makes the system easier to scale and change.

A small example of the layers:

```go
// handler
func (h *CustomerHandler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	customer, err := h.service.CreateCustomer(r.Context(), req.FullName, req.Email, req.Phone)
	if err != nil {
		response.HandleError(w, err)
		return
	}
	response.Success(w, http.StatusCreated, "customer created successfully", customer)
}

// service
func (s *CustomerService) CreateCustomer(ctx context.Context, fullName, email, phone string) (model.Customer, error) {
	customer := model.Customer{FullName: fullName, Email: email, Phone: phone}
	return s.customerRepository.Create(ctx, customer)
}

// repository
func (r *PostgresCustomerRepository) Create(ctx context.Context, customer model.Customer) (model.Customer, error) {
	// INSERT INTO customers ...
}
```

### Repository pattern

The repository layer abstracts database access behind interfaces.

Examples:

- `CustomerRepository`
- `LoanRepository`
- `RepaymentRepository`

Why this matters:

- services are less coupled to Postgres
- data access code is isolated
- it becomes easier to test and refactor

Repository interface example:

```go
type CustomerRepository interface {
	Create(ctx context.Context, customer model.Customer) (model.Customer, error)
	GetByID(ctx context.Context, id int) (model.Customer, error)
	List(ctx context.Context) ([]model.Customer, error)
}
```

Service depending on the interface, not concrete SQL code:

```go
type CustomerService struct {
	customerRepository repository.CustomerRepository
	loanRepository     repository.LoanRepository
}
```

### Observability basics

The middleware layer applies several observability concepts:

- request IDs
- request logging
- panic recovery
- DB-aware health checks

Why this matters:

- requests can be traced through logs
- failures are easier to diagnose
- the service can expose dependency health clearly

Request ID and logging example:

```go
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := fmt.Sprintf("%d", time.Now().UnixNano())
		w.Header().Set("X-Request-ID", requestID)
		ctx := context.WithValue(r.Context(), requestIDKey, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &statusRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(recorder, r)

		log.Printf(
			"request_id=%s method=%s path=%s status=%d duration=%s",
			GetRequestID(r),
			r.Method,
			r.URL.Path,
			recorder.statusCode,
			time.Since(start),
		)
	})
}
```

DB-aware health check example:

```go
func (h *HealthHandler) GetHealth(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := h.db.PingContext(ctx); err != nil {
		response.Error(w, http.StatusServiceUnavailable, "database is unreachable")
		return
	}

	response.Success(w, http.StatusOK, "service is healthy", map[string]string{
		"status":   "ok",
		"database": "up",
	})
}
```

### Defensive validation across layers

This project validates important rules in more than one place:

- service-layer validation for friendly API behavior
- database constraints for persistence-level protection

This is a strong backend/system design pattern because it avoids relying on only one layer for correctness.

Service-layer validation example:

```go
if fullName == "" {
	return model.Customer{}, fmt.Errorf("full_name is required")
}
if email == "" {
	return model.Customer{}, fmt.Errorf("email is required")
}
if phone == "" {
	return model.Customer{}, fmt.Errorf("phone is required")
}
```

Database-layer validation example:

```sql
ALTER TABLE repayments
ADD CONSTRAINT repayments_amount_positive CHECK (amount > 0);
```

### Resource-oriented API design

The endpoints follow resource-oriented API design around:

- customers
- loans
- repayments

Examples:

- `GET /customers/{id}`
- `GET /customers/{id}/loans`
- `GET /loans/{id}`
- `GET /loans/{id}/repayments`
- `GET /repayments/{id}`

This improves API clarity and makes the relationships between entities easier to understand.

Route examples from `main.go`:

```go
mux.HandleFunc("/customers/", func(w http.ResponseWriter, r *http.Request) {
	switch {
	case strings.HasSuffix(r.URL.Path, "/loans"):
		customerHandler.GetCustomerLoans(w, r)
	default:
		customerHandler.GetCustomerByID(w, r)
	}
})

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

## Database constraints

The database currently enforces important integrity rules, including:

- unique customer email
- positive principal amount
- non-negative interest rate
- non-negative total and outstanding amounts
- valid loan status values
- positive repayment amount

This complements service-level validation.

---

## Notes

- migrations are currently run manually using `psql`
- the project uses standard `net/http` instead of an external router
- database constraints are used together with service-level validation
- repayment processing is handled transactionally so repayment creation and loan balance update succeed or fail together

---

## Current project phase

The project is currently in **Phase 2**.

Phase 2 includes:
- Postgres persistence
- migrations
- repository pattern
- richer retrieval endpoints
- middleware
- DB-aware health check
- improved error mapping
- transaction-safe repayment flow

---

## What comes next

Remaining polish before Phase 3:
- final cleanup/refactor
- README refinement if needed

Phase 3 will focus on:
- authentication
- JWT
- role-based access
- stronger testing
- more advanced architecture improvements
```
