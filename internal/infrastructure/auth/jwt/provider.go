// internal/infrastructure/auth/jwt/provider.go
package jwt

import (
	"errors"
	"time"

	"easy-storage/internal/domain/user"

	"github.com/golang-jwt/jwt/v4"
)

// Claims represents the JWT claims
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// Provider handles JWT token generation and validation
type Provider struct {
	secretKey     string
	tokenExpiry   time.Duration
	refreshExpiry time.Duration
}

// NewProvider creates a new JWT provider
func NewProvider(secretKey string, tokenExpiryHours, refreshExpiryDays int) *Provider {
	return &Provider{
		secretKey:     secretKey,
		tokenExpiry:   time.Duration(tokenExpiryHours) * time.Hour,
		refreshExpiry: time.Duration(refreshExpiryDays) * 24 * time.Hour,
	}
}

// GenerateToken generates a new JWT token for a user
func (p *Provider) GenerateToken(user *user.User) (string, error) {
	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(p.tokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(p.secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// GenerateRefreshToken generates a refresh token
func (p *Provider) GenerateRefreshToken(user *user.User) (string, error) {
	claims := &Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(p.refreshExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(p.secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// ValidateToken validates a JWT token
func (p *Provider) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(p.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
