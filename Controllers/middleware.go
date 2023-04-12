package controllers

import (
	"context"
	db "cricHeros/Database"
	u "cricHeros/Utils"
	"net/http"
)

func AdminMiddlerware(originalHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := r.Cookie("token")
		if err != nil {
			u.ShowResponse("Failure", 403, err.Error(), w)
			return
		}

		claims, err := u.DecodeToken(tokenString.Value, w)
		if err != nil {
			u.ShowResponse("Failure", 401, err, w)
			return
		}
		ctx := context.WithValue(r.Context(), "userId", claims.UserId)
		if claims.Role == "admin" {
			originalHandler.ServeHTTP(w, r.WithContext(ctx))
		} else {
			u.ShowResponse("Failed", http.StatusForbidden, "Access denied", w)
		}
	})
}

func LoginMiddlerware(originalHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenString, err := r.Cookie("token")
		var exists bool
		if err != nil {
			u.ShowResponse("Failure", 403, err.Error(), w)
			return
		}
		query := "SELECT EXISTS(SELECT * FROM blacklists where token=?)"
		db.DB.Raw(query, tokenString.Value).Scan(&exists)
		if !exists {
			originalHandler.ServeHTTP(w, r)
		} else {
			u.ShowResponse("Failure", 400, "Token is blacklisted", w)
			return
		}

	})
}
