package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-loan-management-api/internal/auth"
	"go-loan-management-api/internal/model"
)

func TestRequireRoles_RejectsUnauthenticatedRequest(t *testing.T) {
	// This test proves that role middleware does not allow access
	// when there are no auth claims in the request context.
	handler := RequireRoles(model.RoleAdmin)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/customers", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rec.Code)
	}

	if !strings.Contains(rec.Body.String(), "unauthenticated") {
		t.Fatalf("expected unauthenticated error, got %q", rec.Body.String())
	}
}

func TestRequireRoles_RejectsForbiddenRole(t *testing.T) {
	// This test proves that an authenticated user with the wrong role
	// gets a 403 Forbidden response.
	handler := RequireRoles(model.RoleAdmin)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	claims := &auth.Claims{
		UserID: 1,
		Email:  "customer@example.com",
		Role:   model.RoleCustomer,
	}

	req := httptest.NewRequest(http.MethodGet, "/customers", nil)
	req = req.WithContext(auth.WithClaims(req.Context(), claims))

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", rec.Code)
	}

	if !strings.Contains(rec.Body.String(), "forbidden") {
		t.Fatalf("expected forbidden error, got %q", rec.Body.String())
	}
}

func TestRequireRoles_AllowsAllowedRole(t *testing.T) {
	// This test proves that an authenticated user with an allowed role
	// can access the wrapped handler.
	nextCalled := false

	handler := RequireRoles(model.RoleAdmin, model.RoleLoanOfficer)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			nextCalled = true
			w.WriteHeader(http.StatusOK)
		}),
	)

	claims := &auth.Claims{
		UserID: 1,
		Email:  "admin@example.com",
		Role:   model.RoleAdmin,
	}

	req := httptest.NewRequest(http.MethodGet, "/customers", nil)
	req = req.WithContext(auth.WithClaims(req.Context(), claims))

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	if !nextCalled {
		t.Fatalf("expected next handler to be called")
	}
}
