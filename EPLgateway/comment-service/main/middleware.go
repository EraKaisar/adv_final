package main

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func (app *application) requireAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := headerParts[1]
		claims := &jwt.StandardClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(app.config.jwt.secret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", claims.Subject)
		next(w, r.WithContext(ctx))
	}
}

func (app *application) contextGetUserID(r *http.Request) int {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		return 0
	}
	id, _ := strconv.Atoi(userID)
	return id
}
