package auth

import (
	"fmt"
	"go-loan-management-api/internal/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims represents JWT claims used by this API.
type Claims struct {
	UserID     int    `json:"user_id"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	CustomerID *int   `json:"customer_id,omitempty"`
	jwt.RegisteredClaims
}

// GenerateToken creates a signed JWT for a user.
func GenerateToken(user model.User, secret string, expiryHours int) (string, error) {
	now := time.Now()
	expiresAt := now.Add(time.Duration(expiryHours) * time.Hour)

	claims := Claims{
		UserID:     user.ID,
		Email:      user.Email,
		Role:       user.Role,
		CustomerID: user.CustomerID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", user.ID),    // Unique ID for the user
			IssuedAt:  jwt.NewNumericDate(now),       // When it was created
			ExpiresAt: jwt.NewNumericDate(expiresAt), // When it dies
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))

}



// ParseToken validates a signed JWT string and returns its claims.
func ParseToken(tokenString string, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(secret), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
