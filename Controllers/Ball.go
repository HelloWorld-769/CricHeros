package controllers

import (
	db "cricHeros/Database"
	models "cricHeros/Models"
	u "cricHeros/Utils"
	"encoding/json"
	"net/http"
)

func AddBallRecordHandler(scoreCardData models.CardData) {

	var ball models.Balls
	ball.M_ID = scoreCardData.M_ID
	ball.BallType = scoreCardData.Ball_Type
	ball.Runs = scoreCardData.Runs
	ball.P_ID = scoreCardData.Baller

	/*
		//optional
		var matchInning []models.Balls
		db.DB.Where("m_id=?", scoreCardData.M_ID).First(&matchInning)

		if matchInning == nil {
			ball.Inning = "Inning 1"
		}
		if matchInning[0].Inning == "Inning 1" {
			ball.Inning = "Inning 2"
		}

	*/
	var lastBallRecord models.Balls
	query := "SELECT overs FROM balls ORDER BY created_at DESC LIMIT 1"
	db.DB.Raw(query).Scan(&lastBallRecord)
	//fmt.Println("lastBallRecord is ", lastBallRecord)
	if lastBallRecord.Over == 0 {
		ball.BallCount = 1
		ball.Over = 1
		if scoreCardData.Ball_Type == "normal" || scoreCardData.Ball_Type == "wicket" {
			ball.IsValid = "valid"
		}
	} else {
		if scoreCardData.Ball_Type == "normal" || scoreCardData.Ball_Type == "wicket" {
			ball.IsValid = "valid"
			if lastBallRecord.BallCount == 6 {
				ball.Over = lastBallRecord.Over + 1
				ball.BallCount = 1
			} else {
				ball.Over = lastBallRecord.Over
				ball.BallCount = lastBallRecord.BallCount + 1
			}
		} else {
			ball.IsValid = "invalid"
			ball.BallCount = lastBallRecord.BallCount
			ball.Over = lastBallRecord.Over
		}
	}

	db.DB.Create(&ball)
}

// @Description Update the ball
// @Accept json
// @Produce json
// @Tags Ball
// @Param id query string true "Id of the ball"
// @Success 200 {object} models.Response
// @Router /ballUpdate [put]
func UpdateBallRecord(w http.ResponseWriter, r *http.Request) {
	var ballRecord models.Balls
	ball_id := r.URL.Query().Get("id")
	if ball_id == "" {
		u.ShowResponse("Failure", 400, "Please provide ball id", w)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&ballRecord)
	if err != nil {
		u.ShowResponse("Failure", 400, err, w)
		return
	}

	err = db.DB.Where("b_id=?", ball_id).Updates(&ballRecord).Error
	if err != nil {
		u.ShowResponse("Failure", 400, err, w)
		return
	}

	u.ShowResponse("Success", http.StatusOK, ballRecord, w)
}
