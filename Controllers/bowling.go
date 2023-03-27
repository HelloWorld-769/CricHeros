package controllers

import (
	db "cricHeros/Database"
	models "cricHeros/Models"
	"encoding/json"
	"fmt"
	"net/http"
)

// function to add players batting career
func AddBowlingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.URL.Query().Get("id")

	var BowlDetails models.BowlingCareer
	err := json.NewDecoder(r.Body).Decode(&BowlDetails)
	if err != nil {
		fmt.Println("Error in decoding the body")
		return
	}
	BowlDetails.P_ID = id
	db.DB.Create(&BowlDetails)

	json.NewEncoder(w).Encode(&BowlDetails)

}
