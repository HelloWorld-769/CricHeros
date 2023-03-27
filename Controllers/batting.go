package controllers

import (
	db "cricHeros/Database"
	models "cricHeros/Models"
	"encoding/json"
	"fmt"
	"net/http"
)

// function to add players batting career
func AddBattingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.URL.Query().Get("id")

	var BatDetails models.BattingCareer
	err := json.NewDecoder(r.Body).Decode(&BatDetails)
	if err != nil {
		fmt.Println("Error in decoding the body")
		return
	}
	BatDetails.P_ID = id
	db.DB.Create(&BatDetails)

	json.NewEncoder(w).Encode(&BatDetails)

}
