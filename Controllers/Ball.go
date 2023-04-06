package controllers

import (
	db "cricHeros/Database"
	models "cricHeros/Models"
	u "cricHeros/Utils"
	"encoding/json"
	"errors"
	"net/http"

	"gorm.io/gorm"
)

func AddBallRecordHandler(scoreCardData models.CardData) {

	var ball models.Balls
	ball.M_ID = scoreCardData.M_ID
	ball.BallType = scoreCardData.Ball_Type
	ball.Runs = scoreCardData.Runs
	ball.P_ID = scoreCardData.Baller

	var lastBallRecord models.Balls
	err := db.DB.Select("over", "ball_count").Last(&lastBallRecord).Error
	//fmt.Println("lastBallRecord is ", lastBallRecord)
	if errors.Is(err, gorm.ErrRecordNotFound) {

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
// @Success 200 {object} models.Balls
// @Router /ballUpdate [put]
func UpdateBallRecord(w http.ResponseWriter, r *http.Request) {
	ball_id := r.URL.Query().Get("id")
	var ballRecord models.Balls
	err := json.NewDecoder(r.Body).Decode(&ballRecord)
	if err != nil {
		u.ShowResponse("Failure", 400, err.Error(), w)
		return
	}

	err = db.DB.Where("b_id=?", ball_id).Updates(&ballRecord).Error
	if err != nil {
		u.ShowResponse("Failure", http.StatusInternalServerError, err.Error(), w)
		return
	}
	u.ShowResponse("Success", http.StatusOK, ballRecord, w)
}
