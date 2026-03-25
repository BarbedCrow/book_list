package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTProvider struct {
	secret []byte
	ttl    time.Duration
}

func NewJWTProvider(secret string, ttl time.Duration) *JWTProvider {
	return &JWTProvider{secret: []byte(secret), ttl: ttl}
}

func (p *JWTProvider) Generate(userID string) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(p.ttl)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := token.SignedString(p.secret)
	if err != nil {
		return "", fmt.Errorf("jwt generate: %w", err)
	}
	return s, nil
}

func (p *JWTProvider) Validate(tokenStr string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(_ *jwt.Token) (any, error) {
		return p.secret, nil
	})
	if err != nil {
		return "", fmt.Errorf("jwt validate: %w", err)
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return "", errors.New("jwt validate: invalid token")
	}

	return claims.Subject, nil
}
