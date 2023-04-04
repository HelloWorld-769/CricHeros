package controllers

import (
	db "cricHeros/Database"
	models "cricHeros/Models"
	u "cricHeros/Utils"
	"encoding/json"
	"fmt"
	"net/http"
)

// @Description Add player career
// @Accept json
// @Produce json
// @Success 200 {object} models.Career
// @Param id query string true "Player ID"
// @Param playerCareer body models.Career true "Adds Player career"
// @Tags Player
// @Router /addCareer [post]
func AddCareerHandler(w http.ResponseWriter, r *http.Request) {

	EnableCors(&w)
	u.SetHeader(w)
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
