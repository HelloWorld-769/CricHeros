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

func getOvers(player_id string) float64 {
	var count int64 = 0
	db.DB.Model(&models.Balls{}).Where("p_id=?", player_id).Count(&count)
	fmt.Println("count is ", count)
	res := float64(float64(count) / 6)
	return res
}
func ScorecardRecordHandler(w http.ResponseWriter, r *http.Request) {
	var mp = make(map[string]interface{})
	json.NewDecoder(r.Body).Decode(&mp)

	//fmt.Println("map is :", mp)
	//	fmt.Printf("type of runs field is :%T", mp["runs"])

	var bowl models.Balls
	bowl.M_ID = mp["match_id"].(string)
	bowl.BallType = mp["ball_type"].(string)
	bowl.Runs = int64(mp["runs"].(float64))
	bowl.P_ID = mp["bowler"].(string)

	// creating or updating reocrd for bowler
	if val, ok := mp["batsmen"]; ok {
		var existRecord models.ScoreCard
		err := db.DB.Where("p_id=?", val).First(&existRecord).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {

			var record models.ScoreCard
			record.PType = "batsmen"
			record.S_ID = mp["scorecard_id"].(string)
			record.P_ID = val.(string)
			record.RunScored = int64(mp["runs"].(float64))
			if mp["runs"].(float64) == 4 {
				record.Fours += 1
			}
			if mp["runs"].(float64) == 6 {
				record.Sixes += 1
			}
			if mp["ball_type"] == "Normal" {
				record.BPlayed += 1
			}
			record.SR = float64(record.RunScored) / float64(record.BPlayed)

			db.DB.Create(&record)
			u.Encode(w, &record)
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
			existRecord.SR = float64(existRecord.RunScored) / float64(existRecord.BPlayed)
			db.DB.Where("p_id=?", mp["batsmen"]).Updates(&existRecord)
			u.Encode(w, &existRecord)
		}
	} else {
		u.ShowErr("Batsmen not selected", 400, w)
		return
	}

	// creating or updating reocrd for bowler
	if _, ok := mp["bowler"]; ok {
		fmt.Println("Bowler exists and is read to add to database;")
	} else {
		u.ShowErr("Bowler not selected", 400, w)
		return
	}

	bowl.Over = getOvers(mp["bowler"].(string))
	fmt.Println("bowl overs is:", bowl.Over)
	// 	go func() {
	// 		db.DB.Create(&bowl)
	// 	}()
	// }

	db.DB.Create(&bowl)
}
