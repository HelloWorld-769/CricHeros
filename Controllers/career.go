package controllers

import (
	db "cricHeros/Database"
	models "cricHeros/Models"
	"encoding/json"
	"fmt"
	"net/http"
)

// function to add players batting career
func AddCareerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.URL.Query().Get("id")

	var career models.Career
	err := json.NewDecoder(r.Body).Decode(&career)
	if err != nil {
		fmt.Println("Error in decoding the body")
		return
	}
	career.P_ID = id
	db.DB.Create(&career)

	json.NewEncoder(w).Encode(&career)

}
