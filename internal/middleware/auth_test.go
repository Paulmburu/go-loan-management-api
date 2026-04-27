package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-loan-management-api/internal/auth"
	"go-loan-management-api/internal/model"
)

func TestAuthMiddleware_RejectsMissingAuthorizationHeader(t *testing.T) {
	middleware := NewAuthMiddleware("test-secret")

	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	})

	handler := middleware.RequireAuth(next)

	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rec.Code)
	}

	if nextCalled {
		t.Fatalf("expected next handler not to be called")
	}

	body := rec.Body.String()
	if !strings.Contains(body, "missing authorization header") {
		t.Fatalf("expected missing authorization header message, got %q", body)
	}
}

func TestAuthMiddleware_RejectsInvalidBearerFormat(t *testing.T) {
	middleware := NewAuthMiddleware("test-secret")

	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	})

	handler := middleware.RequireAuth(next)

	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	req.Header.Set("Authorization", "Token abc123")

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rec.Code)
	}

	if nextCalled {
		t.Fatalf("expected next handler not to be called")
	}

	body := rec.Body.String()
	if !strings.Contains(body, "invalid authorization header") {
		t.Fatalf("expected invalid authorization header message, got %q", body)
	}
}

func TestAuthMiddleware_RejectsMissingBearerToken(t *testing.T) {
	middleware := NewAuthMiddleware("test-secret")

	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	})

	handler := middleware.RequireAuth(next)

	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	req.Header.Set("Authorization", "Bearer ")

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rec.Code)
	}

	if nextCalled {
		t.Fatalf("expected next handler not to be called")
	}

	body := rec.Body.String()
	if !strings.Contains(body, "missing bearer token") {
		t.Fatalf("expected missing bearer token message, got %q", body)
	}
}

func TestAuthMiddleware_RejectsInvalidToken(t *testing.T) {
	middleware := NewAuthMiddleware("test-secret")

	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	})

	handler := middleware.RequireAuth(next)

	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rec.Code)
	}

	if nextCalled {
		t.Fatalf("expected next handler not to be called")
	}

	body := rec.Body.String()
	if !strings.Contains(body, "invalid or expired token") {
		t.Fatalf("expected invalid or expired token message, got %q", body)
	}
}

func TestAuthMiddleware_AllowsValidToken(t *testing.T) {
	user := model.User{
		ID:       1,
		FullName: "Paul Mburu",
		Email:    "paul@example.com",
		Role:     model.RoleAdmin,
	}

	token, err := auth.GenerateToken(user, "test-secret", 24)
	if err != nil {
		t.Fatalf("failed to generate test token: %v", err)
	}

	middleware := NewAuthMiddleware("test-secret")

	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true

		claims, ok := auth.GetClaims(r.Context())
		if !ok {
			t.Fatalf("expected claims to be present in request context")
		}

		if claims.UserID != 1 {
			t.Fatalf("expected user ID 1, got %d", claims.UserID)
		}

		if claims.Email != "paul@example.com" {
			t.Fatalf("expected email paul@example.com, got %s", claims.Email)
		}

		if claims.Role != model.RoleAdmin {
			t.Fatalf("expected role %s, got %s", model.RoleAdmin, claims.Role)
		}

		w.WriteHeader(http.StatusOK)
	})

	handler := middleware.RequireAuth(next)

	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	if !nextCalled {
		t.Fatalf("expected next handler to be called")
	}
}
