package controllers

import (
	db "cricHeros/Database"
	models "cricHeros/Models"
	"encoding/json"
	"net/http"
)

// not completed
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var cred models.Credentials
	json.NewDecoder(r.Body).Decode(&cred)

	db.DB.Create(&cred)
	json.NewEncoder(w).Encode(&cred)
}
