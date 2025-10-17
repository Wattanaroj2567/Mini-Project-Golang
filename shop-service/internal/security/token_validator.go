package security

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// Claims defines the JWT payload we expect from users-service.
type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// Validator validates JWT tokens using a shared secret.
type Validator struct {
	secret []byte
}

// NewValidator constructs a Validator with the provided secret.
func NewValidator(secret string) (*Validator, error) {
	if secret == "" {
		return nil, fmt.Errorf("jwt secret must not be empty")
	}
	return &Validator{secret: []byte(secret)}, nil
}

// Validate parses and verifies the supplied JWT string and returns the claims.
func (v *Validator) Validate(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return v.secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("token is not valid")
	}

	return claims, nil
}
