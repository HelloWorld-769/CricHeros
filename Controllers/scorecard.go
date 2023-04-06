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

// @Description stores players info in scorecard
// @Accept json
// @Success 200 {object} models.Response
// @Param details body models.CardData true "ScoreCard details"
// @Tags Scorecard
// @Router /addToScoreCard [post]
func ScorecardRecordHandler(w http.ResponseWriter, r *http.Request) {
	u.SetHeader(w)
	u.EnableCors(&w)
	//var mp = make(map[string]interface{})
	var scoreCardData models.CardData
	err := json.NewDecoder(r.Body).Decode(&scoreCardData)
	if err != nil {
		u.ShowResponse("Failure", 400, err.Error(), w)
		return
	}

	validationErr := u.CheckValidation(scoreCardData)
	if validationErr != nil {
		u.ShowResponse("Failure", 400, validationErr.Error(), w)
		return
	}
	var matchMapping models.MatchRecord
	err = db.DB.Where("m_id=?", scoreCardData.M_ID).First(&matchMapping).Error
	if err != nil {
		u.ShowResponse("Failure", 400, err.Error(), w)
		return
	}
	go AddBallRecordHandler(scoreCardData)
	// creating or updating reocrd for bowler
	if scoreCardData.Batsmen != "" {
		var existRecord models.ScoreCard
		err := db.DB.Where("p_id=?", scoreCardData.Batsmen).First(&existRecord).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {

			var batsmenRecord models.ScoreCard
			batsmenRecord.PType = "batsmen"
			batsmenRecord.S_ID = matchMapping.S_ID
			batsmenRecord.P_ID = scoreCardData.Batsmen
			batsmenRecord.RunScored = scoreCardData.Runs
			if scoreCardData.Runs == 4 {
				batsmenRecord.Fours += 1
			}
			if scoreCardData.Runs == 6 {
				batsmenRecord.Sixes += 1
			}
			if scoreCardData.Ball_Type == "normal" {
				batsmenRecord.BPlayed += 1
			}
			if scoreCardData.Ball_Type == "wicket" {
				batsmenRecord.IsOut = "Out"
			}
			batsmenRecord.SR = u.RoundFloat((float64(batsmenRecord.RunScored)/float64(batsmenRecord.BPlayed))*100, 3)
			err = db.DB.Create(&batsmenRecord).Error
			if err != nil {
				u.ShowResponse("Failure", 400, err.Error(), w)
				return
			}
			u.ShowResponse("Success", 200, batsmenRecord, w)
		} else {
			//update the scorecard for that user
			existRecord.RunScored = existRecord.RunScored - scoreCardData.PrevRuns + scoreCardData.Runs
			if scoreCardData.Runs == 4 {
				existRecord.Fours += 1
			}
			if scoreCardData.Runs == 6 {
				existRecord.Sixes += 1
			}
			if scoreCardData.Ball_Type == "normal" {
				existRecord.BPlayed += 1
			}
			if scoreCardData.Ball_Type == "wicket" {
				existRecord.IsOut = "Out"
			}
			existRecord.SR = u.RoundFloat(float64(existRecord.RunScored)/float64(existRecord.BPlayed), 3)
			err = db.DB.Where("p_id=?", scoreCardData.Batsmen).Updates(&existRecord).Error
			if err != nil {
				u.ShowResponse("Failure", 400, err.Error(), w)
				return
			}
			u.ShowResponse("Success", 200, existRecord, w)
		}
	} else {
		u.ShowResponse("Failure", 400, "Batsmen not selected", w)
		return
	}

	// creating or updating record for bowler
	if scoreCardData.Baller != "" {
		var existRecord models.ScoreCard
		err := db.DB.Where("p_id=?", scoreCardData.Baller).First(&existRecord).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			var bowlerRecord models.ScoreCard
			bowlerRecord.PType = "bowler"
			bowlerRecord.S_ID = matchMapping.S_ID
			bowlerRecord.P_ID = scoreCardData.Baller
			bowlerRecord.RunGiven = scoreCardData.Runs
			if scoreCardData.Ball_Type == "no_ball" {
				bowlerRecord.NB += 1
			} else if scoreCardData.Ball_Type == "wicket" {
				bowlerRecord.Wickets += 1
			} else if scoreCardData.Ball_Type == "wide_ball" {
				bowlerRecord.WD += 1
			}
			err = db.DB.Create(&bowlerRecord).Error
			if err != nil {
				u.ShowResponse("Failure", 400, err.Error(), w)
				return
			}
			u.ShowResponse("Success", 200, bowlerRecord, w)
		} else {
			existRecord.RunGiven += scoreCardData.Runs
			if scoreCardData.Ball_Type == "no_ball" {
				existRecord.NB += 1
			} else if scoreCardData.Ball_Type == "wicket" {
				existRecord.Wickets += 1
			} else if scoreCardData.Ball_Type == "wide_ball" {
				existRecord.WD += 1
			}
			if getOvers(scoreCardData.Baller) != 0 {
				existRecord.OBowled = getOvers(scoreCardData.Baller)
			}
			if getMaidenOvers(scoreCardData.Baller) != 0 {
				existRecord.MOvers = getMaidenOvers(scoreCardData.Baller)
			}
			existRecord.Eco = float64(existRecord.RunGiven) / float64(existRecord.OBowled)
			err = db.DB.Where("p_id=?", scoreCardData.Baller).Updates(&existRecord).Error
			if err != nil {
				u.ShowResponse("Failure", 400, err.Error(), w)
				return
			}
			u.ShowResponse("Success", 200, existRecord, w)
		}
	} else {
		u.ShowResponse("Failure", 400, "Bowler not selected", w)
		return
	}

}

// @Description Shows the score card for the current matcha
// @Accept json
// @Success 200 {object} models.Response
// @Param match_id body object true "Id of the match whose scoredcard is to be viewed"
// @Tags Scorecard
// @Router /showScoreCard [post]
func ShowScoreCardHandler(w http.ResponseWriter, r *http.Request) {
	u.SetHeader(w)
	u.EnableCors(&w)
	var mp = make(map[string]string)
	err := json.NewDecoder(r.Body).Decode(&mp)
	if err != nil {
		u.ShowResponse("Failure", 400, err.Error(), w)
		return
	}

	var matchMapping models.MatchRecord
	err = db.DB.Where("m_id=?", mp["matchId"]).First(&matchMapping).Error
	if err != nil {
		u.ShowResponse("Failure", 400, err.Error(), w)
		return
	}

	var matchScoreRecord []models.ScoreCard
	err = db.DB.Where("s_id=?", matchMapping.S_ID).Find(&matchScoreRecord).Error
	if err != nil {
		u.ShowResponse("Failure", 400, err.Error(), w)
		return
	}

	u.ShowResponse("Success", 200, matchScoreRecord, w)
}
