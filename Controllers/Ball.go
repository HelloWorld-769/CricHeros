package controllers

import (
	db "cricHeros/Database"
	models "cricHeros/Models"
	"errors"

	"gorm.io/gorm"
)

func AddBallRecordHandler(mp map[string]interface{}) {
	// var mp = make(map[string]interface{})
	// json.NewDecoder(r.Body).Decode(&mp)
	//fmt.Println("map is :", mp)
	var ball models.Balls
	ball.M_ID = mp["match_id"].(string)
	ball.BallType = mp["ball_type"].(string)
	ball.Runs = int64(mp["runs"].(float64))
	ball.P_ID = mp["bowler"].(string)

	var lastBallRecord models.Balls
	err := db.DB.Select("over", "ball_count").Last(&lastBallRecord).Error
	//fmt.Println("lastBallRecord is ", lastBallRecord)
	if errors.Is(err, gorm.ErrRecordNotFound) {

		ball.BallCount = 1
		ball.Over = 1
		if mp["ball_type"].(string) == "normal" || mp["ball_type"].(string) == "wicket" {
			ball.IsValid = "valid"
		}
	} else {
		if mp["ball_type"].(string) == "normal" || mp["ball_type"].(string) == "wicket" {
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
