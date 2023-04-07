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
// @Success 200 {object} models.Response
// @Param id query string true "Player ID"
// @Param playerCareer body models.Career true "Adds Player career"
// @Tags Player
// @Router /addCareer [post]
func AddCareerHandler(w http.ResponseWriter, r *http.Request) {

	u.EnableCors(&w)
	u.SetHeader(w)
	var career models.Career
	id := r.URL.Query().Get("id")
	if id == "" {
		u.ShowResponse("Failure", 400, "Please provide id", w)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&career)
	if err != nil {
		u.ShowResponse("Failure", 400, err, w)
		return
	}

	career.P_ID = id
	err = db.DB.Create(&career).Error
	if err != nil {
		u.ShowResponse("Failure", 500, "Internal Server error", w)
		return
	}

	u.ShowResponse("Success", http.StatusOK, &career, w)

}
