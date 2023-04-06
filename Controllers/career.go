package controllers

import (
	db "cricHeros/Database"
	models "cricHeros/Models"
	u "cricHeros/Utils"
	"encoding/json"
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

	u.EnableCors(&w)
	u.SetHeader(w)
	id := r.URL.Query().Get("id")

	var career models.Career
	err := json.NewDecoder(r.Body).Decode(&career)
	if err != nil {
		u.ShowResponse("Failure", 400, err.Error(), w)
		return
	}
	career.P_ID = id
	err = db.DB.Create(&career).Error
	if err != nil {
		u.ShowResponse("Failure", 400, err.Error(), w)
		return
	}

	u.ShowResponse("Success", http.StatusOK, &career, w)

}
