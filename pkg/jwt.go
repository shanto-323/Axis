package pkg

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/shanto-323/axis/config"
)

type CustomClaim struct {
	ID uuid.UUID `json:"id"`
	jwt.RegisteredClaims
}

func CreateAccessToken(cfg *config.Config, id uuid.UUID) (string, error) {
	claims := CustomClaim{
		ID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    cfg.Observability.ServiceName,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Server.JwtKey))
}

func ValidateToken(cfg *config.Config, tokenString string) (*CustomClaim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaim{}, func(token *jwt.Token) (any, error) {

		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method")
		}

		return []byte(cfg.Server.JwtKey), nil
	})

	if err != nil {
		return nil, err
	} else if claims, ok := token.Claims.(*CustomClaim); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("invalid token")
	}
}
