package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func generateJWT(claims Claims, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":      claims.UserID,
		"email":        claims.Email,
		"role":         claims.Role,
		"reference_id": claims.ReferenceID,
		"exp":          time.Now().Add(24 * time.Hour).Unix(),
		"iat":          time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func validateJWT(tokenString, secret string) (*Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid user_id in token")
	}

	referenceIDFloat, _ := claims["reference_id"].(float64)

	return &Claims{
		UserID:      int64(userIDFloat),
		Email:       claims["email"].(string),
		Role:        claims["role"].(string),
		ReferenceID: int64(referenceIDFloat),
	}, nil
}
