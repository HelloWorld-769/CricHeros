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

// @Description Create the match between the teams
// @Accept json
// @Produces json
// @Success 200 {object} models.Match
// @Param match body models.Match true "Match details"
// @Tags Match
// @Router /createMatch [post]
func CreateMatchHandler(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
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
	db.DB.Create(&matchRecord)

	u.Encode(w, &match)
}

// @Description Show the list of matches
// @Accept json
// @Produces json
// @Success 200 {object} models.Match
// @Tags Match
// @Router /showMatch [post]
func ShowMatchHandler(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	u.SetHeader(w)
	var matches []models.Match
	db.DB.Find(&matches)
	u.Encode(w, &matches)
}

// @Description Ends the match and updates the scorecard of every player
// @Accept json
// @Success 200 {object} models.Match
// @Param match_id body object true "Id of the match to end"
// @Tags Match
// @Router /endMatch [post]
func EndMatchHandler(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	u.SetHeader(w)
	var mp = make(map[string]string)
	json.NewDecoder(r.Body).Decode(&mp)
	var matchData models.Match
	db.DB.Where("s_id", mp["match_id"]).Find(&matchData)
	matchData.Status = "Completed"
	var scorecard models.MatchRecord
	db.DB.Where("m_id=?", mp["match_id"]).First(&scorecard)
	var records []models.ScoreCard
	db.DB.Where("s_id", scorecard.S_ID).Find(&records)
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
