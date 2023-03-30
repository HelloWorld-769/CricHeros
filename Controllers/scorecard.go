package controllers

import (
	db "cricHeros/Database"
	models "cricHeros/Models"
	u "cricHeros/Utils"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

func getOvers(player_id string) int64 {
	var count int64 = 0
	db.DB.Model(&models.Balls{}).Where("p_id=? AND is_valid=?", player_id, "valid").Count(&count)
	fmt.Println("overs balls is ", count)
	if count%6 == 0 {
		return count / 6
	}
	return 0
}

func getMaidenOvers(player_id string) int64 {
	var count int64
	db.DB.Model(&models.Balls{}).Where("p_id=? AND is_valid=? AND runs=?", player_id, "valid", 0).Count(&count)
	fmt.Println("Maiden balls count: ", count)
	if count%6 == 0 {
		return count / 6
	}
	return 0
}
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
func ScorecardRecordHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var mp = make(map[string]interface{})
	json.NewDecoder(r.Body).Decode(&mp)

	AddBallRecordHandler(mp)
	// creating or updating reocrd for bowler
	if val, ok := mp["batsmen"].(string); ok {
		var existRecord models.ScoreCard
		err := db.DB.Where("p_id=?", val).First(&existRecord).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {

			var batsmenRecord models.ScoreCard
			batsmenRecord.PType = "batsmen"
			batsmenRecord.S_ID = mp["scorecard_id"].(string)
			batsmenRecord.P_ID = val
			batsmenRecord.RunScored = int64(mp["runs"].(float64))
			if mp["runs"].(float64) == 4 {
				batsmenRecord.Fours += 1
			}
			if mp["runs"].(float64) == 6 {
				batsmenRecord.Sixes += 1
			}
			if mp["ball_type"] == "normal" {
				batsmenRecord.BPlayed += 1
			}
			if mp["ball_type"] == "wicket" {
				batsmenRecord.IsOut = "Out"
			}
			batsmenRecord.SR = u.RoundFloat((float64(batsmenRecord.RunScored)/float64(batsmenRecord.BPlayed))*100, 3)
			db.DB.Create(&batsmenRecord)
			u.Encode(w, &batsmenRecord)
		} else {
			//update the scorecard for that user
			existRecord.RunScored += int64(mp["runs"].(float64))
			if mp["runs"].(float64) == 4 {
				existRecord.Fours += 1
			}
			if mp["runs"].(float64) == 6 {
				existRecord.Sixes += 1
			}
			if mp["ball_type"] == "normal" {
				existRecord.BPlayed += 1
			}
			if mp["ball_type"] == "wicket" {
				existRecord.IsOut = "Out"
			}
			existRecord.SR = u.RoundFloat(float64(existRecord.RunScored)/float64(existRecord.BPlayed), 3)
			db.DB.Where("p_id=?", mp["batsmen"]).Updates(&existRecord)
			u.Encode(w, &existRecord)
		}
	} else {
		u.ShowErr("Batsmen not selected", 400, w)
		return
	}

	// creating or updating record for bowler
	if val, ok := mp["bowler"].(string); ok {
		var existRecord models.ScoreCard
		err := db.DB.Where("p_id=?", val).First(&existRecord).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			var bowlerRecord models.ScoreCard
			bowlerRecord.PType = "bowler"
			bowlerRecord.S_ID = mp["scorecard_id"].(string)
			bowlerRecord.P_ID = val
			bowlerRecord.RunGiven = int64(mp["runs"].(float64))
			if mp["ball_type"] == "no_ball" {
				bowlerRecord.NB += 1
			} else if mp["ball_type"] == "wicket" {
				bowlerRecord.Wickets += 1
			} else if mp["ball_type"] == "wide_ball" {
				bowlerRecord.WD += 1
			}
			db.DB.Create(&bowlerRecord)
		} else {
			existRecord.RunGiven += int64(mp["runs"].(float64))
			if mp["ball_type"] == "no_ball" {
				existRecord.NB += 1
			} else if mp["ball_type"] == "wicket" {
				existRecord.Wickets += 1
			} else if mp["ball_type"] == "wide_ball" {
				existRecord.WD += 1
			}
			if getOvers(mp["bowler"].(string)) != 0 {
				existRecord.OBowled = getOvers(mp["bowler"].(string))
			}
			if getMaidenOvers(mp["bowler"].(string)) != 0 {
				existRecord.MOvers = getMaidenOvers(mp["bowler"].(string))
			}
			//existRecord.Eco=float64(existRecord.RunGiven)/float64(existRecord.OBowled)
			db.DB.Where("p_id=?", mp["bowler"].(string)).Updates(&existRecord)
			fmt.Fprintln(w, "Baller record is :")
			u.Encode(w, &existRecord)
		}
	} else {
		u.ShowErr("Bowler not selected", 400, w)
		return
	}

}

func ShowScoreCard(w http.ResponseWriter, r *http.Request) {

}
