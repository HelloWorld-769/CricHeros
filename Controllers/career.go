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
// @Param playerCareer body models.Career true "Adds Player career"
// @Tags Player
// @Router /addCareer [post]
func AddCareerHandler(w http.ResponseWriter, r *http.Request) {

	u.EnableCors(&w)
	u.SetHeader(w)
	var career models.Career
	var exists bool
	err := json.NewDecoder(r.Body).Decode(&career)
	if err != nil {
		u.ShowResponse("Failure", 400, err, w)
		return
	}
	query := "SELECT EXISTS(SELECT * FROM careers where p_id=?)"
	db.DB.Raw(query, career.P_ID).Scan(&exists)
	if exists {
		err = db.DB.Where("p_id=?", career.P_ID).Updates(&career).Error
		if err != nil {
			u.ShowResponse("Failure", 400, err, w)
			return
		}

	} else {
		err = db.DB.Create(&career).Error
		if err != nil {
			u.ShowResponse("Failure", 500, "Internal Server error", w)
			return
		}
	}

	u.ShowResponse("Success", http.StatusOK, &career, w)

}
