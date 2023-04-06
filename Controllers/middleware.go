package controllers

import (
	"context"
	models "cricHeros/Models"
	u "cricHeros/Utils"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v4"
)

func AdminMiddlerware(originalHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("token")
		claims := &models.Claims{}

		parsedToken, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("error")
			}
			return []byte(os.Getenv("SECRET_KEY")), nil
		})
		if err != nil || !parsedToken.Valid {
			http.Error(w, "Invalid or expired token", http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), "claims", &claims)
		if claims.Role == "admin" {
			originalHandler.ServeHTTP(w, r.WithContext(ctx))
		} else {
			u.ShowResponse("Failed", http.StatusForbidden, "Access denied", w)
		}
	})
}
