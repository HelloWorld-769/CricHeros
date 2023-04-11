package controllers

import (
	"context"
	db "cricHeros/Database"
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
			u.ShowResponse("Failure", 401, "Invalid Token", w)
			return
		}
		ctx := context.WithValue(r.Context(), "userId", claims.UserID)
		if claims.Role == "admin" {
			originalHandler.ServeHTTP(w, r.WithContext(ctx))
		} else {
			u.ShowResponse("Failed", http.StatusForbidden, "Access denied", w)
		}
	})
}

func LoginMiddlerware(originalHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("token")
		claims := &models.Claims{}
		var exists bool
		parsedToken, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("error")
			}
			return []byte(os.Getenv("SECRET_KEY")), nil
		})
		if err != nil || !parsedToken.Valid {
			u.ShowResponse("Failure", 401, "Invalid Token", w)
			return
		}
		var cred models.Credential
		err = db.DB.Where("user_id=?", claims.UserID).First(&cred).Error
		if err != nil {
			u.ShowResponse("Failure", 400, err, w)
			return
		}

		query := "SELECT EXISTS(SELECT * FROM blacklists where token=?)"
		db.DB.Raw(query, tokenString).Scan(&exists)
		if !exists {
			originalHandler.ServeHTTP(w, r)
		} else {
			u.ShowResponse("Failure", 400, "Token is blacklisted", w)
			return
		}

	})
}
