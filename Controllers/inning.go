package controllers

import (
	db "cricHeros/Database"
	models "cricHeros/Models"
	u "cricHeros/Utils"
	"encoding/json"
	"net/http"
)

func EndInningHandler(w http.ResponseWriter, r *http.Request) {
	u.SetHeader(w)
	var mp = make(map[string]string)
	json.NewDecoder(r.Body).Decode(&mp)

	var teamData []models.Team
	db.DB.Where("t_id=?", mp["team_id"]).Find(&teamData)
	totalScore := 0
	for _, player := range teamData {
		var playerRecord models.ScoreCard
		db.DB.Where("p_id=?", player.P_ID).Find(&playerRecord)
		totalScore += int(playerRecord.RunScored)
	}
	inning := models.Inning{
		T_ID:   mp["team_id"],
		TScore: int64(totalScore),
	}

	db.DB.Create(&inning)
}
