package controllers

import (
	db "cricHeros/Database"
	models "cricHeros/Models"
	u "cricHeros/Utils"
	"encoding/json"
	"net/http"
)

// @Description Ends the current team innings
// @Accept json
// @Param team_id body object true "Id of the team to end its inning"
// @Tags Inning
// @Success 200
// @Router /endInning [post]
func EndInningHandler(w http.ResponseWriter, r *http.Request) {
	u.EnableCors(&w)
	u.SetHeader(w)
	var mp = make(map[string]string)

	err := json.NewDecoder(r.Body).Decode(&mp)
	if err != nil {
		u.ShowResponse("Failure", 400, err.Error(), w)
		return
	}

	var teamData []models.Team
	err = db.DB.Where("t_id=?", mp["teamId"]).Find(&teamData).Error
	if err != nil {
		u.ShowResponse("Failure", 400, err.Error(), w)
		return
	}
	totalScore := 0
	for _, player := range teamData {
		var playerRecord models.ScoreCard
		db.DB.Where("p_id=?", player.P_ID).Find(&playerRecord)
		totalScore += int(playerRecord.RunScored)
	}
	inning := models.Inning{
		M_ID:   mp["matchId"],
		T_ID:   mp["teamId"],
		TScore: int64(totalScore),
	}

	err = db.DB.Create(&inning).Error
	if err != nil {
		u.ShowResponse("Failure", 400, err.Error(), w)
		return
	}

	u.ShowResponse("Success", http.StatusOK, inning, w)
}

func EndInningHandler2(mp map[string]string, w http.ResponseWriter) {

	var teamData []models.Team
	err := db.DB.Where("t_id=?", mp["teamId"]).Find(&teamData).Error
	if err != nil {
		u.ShowResponse("Failure", 400, err.Error(), w)
		return
	}
	totalScore := 0
	for _, player := range teamData {
		var playerRecord models.ScoreCard
		db.DB.Where("p_id=?", player.P_ID).Find(&playerRecord)
		totalScore += int(playerRecord.RunScored)
	}
	inning := models.Inning{
		M_ID:   mp["matchId"],
		T_ID:   mp["teamId"],
		TScore: int64(totalScore),
	}

	err = db.DB.Create(&inning).Error
	if err != nil {
		u.ShowResponse("Failure", 400, err.Error(), w)
		return
	}

	u.ShowResponse("Success", http.StatusOK, inning, w)
}
