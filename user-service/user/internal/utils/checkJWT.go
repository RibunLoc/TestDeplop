package utils

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func ExtractToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("missing Authorization header")
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", errors.New("invalid Authorization format")
	}
	return parts[1], nil
}

func VerifyJWT(r *http.Request) (jwt.MapClaims, error) {
	tokenStr, err := ExtractToken(r)
	if err != nil {
		return nil, err
	}

	envPath, err := LoadEnv()
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

    // Load the .env file
    err = godotenv.Load(envPath)
	if err != nil {
		log.Fatalf("Error loading .env file from %s: %v", envPath, err)
	}
	
	jwtSecret := []byte(os.Getenv("JWT_SECRET_KEY"))
	if jwtSecret == nil {
		return nil, errors.New("JWT secret key not found in environment variables")
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})
	
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("cannot parse claims")
	}

	return claims, nil
}
