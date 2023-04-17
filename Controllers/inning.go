package controllers

import (
	db "cricHeros/Database"
	models "cricHeros/Models"
	u "cricHeros/Utils"
	"encoding/json"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// @Description Ends the current team innings
// @Accept json
// @Param matchDetils body string true "Id of the team to end its inning" SchemaExample({\n "matchId":"string",\n "teamId":"string"\n})
// @Tags Inning
// @Success 200 {object} models.Response
// @Router /endInning [post]
func EndInningHandler(w http.ResponseWriter, r *http.Request) {
	u.EnableCors(&w)
	u.SetHeader(w)
	var mp = make(map[string]interface{})
	var teamData []models.Team
	var exists bool
	err := json.NewDecoder(r.Body).Decode(&mp)
	if err != nil {
		u.ShowResponse("Failure", 400, err, w)
		return
	}
	err = validation.Validate(mp,
		validation.Map(
			validation.Key("teamId", validation.Required),
			validation.Key("matchId", validation.Required),
		),
	)
	if err != nil {
		u.ShowResponse("Failure", 400, err, w)
		return
	}

	query := "SELECT EXISTS(SELECT * FROM matches WHERE m_id=? AND (t1_id=? OR t2_id='?))"
	db.DB.Raw(query, mp["matchId"], mp["teamId"], mp["teamId"]).Scan(&exists)

	if !exists {
		u.ShowResponse("Failure", 400, "The given team id is not linked with the match..", w)
		return

	}

	err = db.DB.Where("t_id=?", mp["teamId"].(string)).Find(&teamData).Error
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
		M_ID:   mp["matchId"].(string),
		T_ID:   mp["teamId"].(string),
		TScore: int64(totalScore),
	}

	err = db.DB.Create(&inning).Error
	if err != nil {
		u.ShowResponse("Failure", 400, err, w)
		return
	}
	u.ShowResponse("Success", http.StatusOK, inning, w)
}

func EndInningHandler2(mp map[string]interface{}, w http.ResponseWriter) {

	var teamData []models.Team
	err := db.DB.Where("t_id=?", mp["teamId"].(string)).Find(&teamData).Error
	if err != nil {
		u.ShowResponse("Failure", 400, err, w)
		return
	}

	totalScore := 0
	for _, player := range teamData {
		var playerRecord models.ScoreCard
		db.DB.Where("p_id=?", player.P_ID).Find(&playerRecord)
		totalScore += int(playerRecord.RunScored)
	}
	inning := models.Inning{
		M_ID:   mp["matchId"].(string),
		T_ID:   mp["teamId"].(string),
		TScore: int64(totalScore),
	}

	err = db.DB.Create(&inning).Error
	if err != nil {
		u.ShowResponse("Failure", 500, "Internal server error", w)
		return
	}

	//Ask from sir

	var ball models.Balls
	ball.Over = 0
	ball.BallCount = 0

	db.DB.Create(&ball)

	//u.ShowResponse("Success", http.StatusOK, inning, w)
}
