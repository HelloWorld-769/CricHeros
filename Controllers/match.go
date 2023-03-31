package controllers

import (
	db "cricHeros/Database"
	models "cricHeros/Models"
	u "cricHeros/Utils"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"
)

func CreateMatchHandler(w http.ResponseWriter, r *http.Request) {
	u.SetHeader(w)
	var match models.Match
	json.NewDecoder(r.Body).Decode(&match)
	var count int64
	query := "SELECT * FROM teams where t_id=? and p_id IN (SELECT p_id FROM teams WHERE t_id=? )"
	db.DB.Raw(query, match.T1_ID, match.T2_ID).Count(&count)
	fmt.Println("Count isL: ", count)
	if count != 0 {
		u.ShowErr("Common record found Please Edit your teams", 400, w)
		return
	}
	now := time.Now()
	match.Date = now.Format("02 Jan 2006")
	db.DB.Create(&match)
	matchRecord := models.MatchRecord{
		M_ID: match.M_ID,
		S_ID: match.S_ID,
	}
	// tossRecord := models.Toss{
	// 	M_ID:  match.M_ID,
	// 	T1_ID: match.T1_ID,
	// 	T2_ID: match.T2_ID,
	// }
	db.DB.Create(&matchRecord)
	//db.DB.Create(&tossRecord)

	u.Encode(w, &match)
}
func ShowMatchHandler(w http.ResponseWriter, r *http.Request) {
	u.SetHeader(w)
	var matches []models.Match
	db.DB.Find(&matches)
	u.Encode(w, &matches)
}

func EndMatchHandler(w http.ResponseWriter, r *http.Request) {
	u.SetHeader(w)
	var mp = make(map[string]string)
	json.NewDecoder(r.Body).Decode(&mp)
	var matchData models.Match
	db.DB.Where("s_id", mp["match_id"]).Find(&matchData)
	matchData.Status = "Completed"
	var records []models.ScoreCard
	db.DB.Where("s_id", mp["scorecard_id"]).Find(&records)
	for _, record := range records {
		//fmt.Println("Player id is:", record.P_ID)
		var pCareer models.Career
		db.DB.Where("p_id=?", record.P_ID).First(&pCareer)
		pCareer.MPlayed += 1
		if record.PType == "batsmen" {
			pCareer.RunScored += record.RunScored
			pCareer.HScored = int64(math.Max(float64(pCareer.HScored), float64(record.RunScored)))
			pCareer.BallsFaced += record.BPlayed
			if record.RunScored >= 50 && record.RunScored < 100 {
				pCareer.Fifites += 1
			} else if record.RunScored >= 100 && record.RunScored < 200 {
				pCareer.Hundreds += 1
			} else if record.RunScored >= 200 {
				pCareer.TwoHundreds += 1
			}
			pCareer.AvgScore = u.RoundFloat(float64(pCareer.RunScored)/float64(pCareer.MPlayed), 2)

		} else if record.PType == "bowler" {
			pCareer.RConced += record.RunGiven
			pCareer.Wickets += record.Wickets
			pCareer.BBowl += record.OBowled * 6
			pCareer.BowlAvg = u.RoundFloat(float64(pCareer.RConced)/float64(pCareer.Wickets), 2)
			pCareer.Economy = u.RoundFloat(float64(pCareer.RConced)/float64(pCareer.BBowl/6), 2)
		}
		db.DB.Where("p_id=?", record.P_ID).Updates(&pCareer)
	}

	//fmt.Println("records in scorecard is:", records)

}
