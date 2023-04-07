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

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// @Description Create the match between the teams
// @Accept json
// @Produces json
// @Success 200 {object} models.Response
// @Param match body models.Match true "Match details"
// @Tags Match
// @Router /createMatch [post]
func CreateMatchHandler(w http.ResponseWriter, r *http.Request) {
	u.EnableCors(&w)
	u.SetHeader(w)
	claims := r.Context().Value("claims").(*models.Claims)
	var match models.Match
	err := json.NewDecoder(r.Body).Decode(&match)
	if err != nil {
		u.ShowResponse("Failure", 400, err, w)
		return
	}
	var count int64
	query := "SELECT * FROM teams where t_id=? and p_id IN (SELECT p_id FROM teams WHERE t_id=? )"
	db.DB.Raw(query, match.T1_ID, match.T2_ID).Count(&count)
	fmt.Println("Count is: ", count)
	if count != 0 {
		u.ShowResponse("Failure", 400, "Common record found Please Edit your teams", w)
		return
	}
	match.U_ID = claims.UserID
	now := time.Now()
	match.Date = now.Format("02 Jan 2006")
	err = db.DB.Create(&match).Error
	if err != nil {
		u.ShowResponse("Failure", 400, err, w)
		return
	}
	matchRecord := models.MatchRecord{
		M_ID: match.M_ID,
		S_ID: match.S_ID,
	}
	err = db.DB.Create(&matchRecord).Error
	if err != nil {
		u.ShowResponse("Failure", 500, "Inteternal Server Error", w)
		return
	}
	u.ShowResponse("Success", http.StatusOK, match, w)
}

// @Description Show the list of matches
// @Accept json
// @Produces json
// @Success 200 {object} models.Response
// @Tags Match
// @Router /showMatch [post]
func ShowMatchHandler(w http.ResponseWriter, r *http.Request) {
	u.EnableCors(&w)
	u.SetHeader(w)
	var matches []models.Match
	err := db.DB.Find(&matches).Error
	if err != nil {
		u.ShowResponse("Failure", 400, err, w)
		return
	}
	u.ShowResponse("Success", http.StatusOK, matches, w)
}

// @Description Ends the match and updates the scorecard of every player
// @Accept json
// @Success 200 {object} models.Response
// @Param match_id body object true "Id of the match to end"
// @Tags Match
// @Router /endMatch [post]
func EndMatchHandler(w http.ResponseWriter, r *http.Request) {
	u.EnableCors(&w)
	u.SetHeader(w)
	claims := r.Context().Value("claims").(*models.Claims)
	var mp = make(map[string]interface{})
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
	EndInningHandler2(mp, w)

	//get match data to update its status to completed
	var matchData models.Match
	err = db.DB.Where("s_id", mp["matchId"].(string)).Find(&matchData).Error
	if err != nil {
		u.ShowResponse("Failure", 400, err, w)
		return
	}

	err = validation.Validate(mp, validation.Map(validation.Key("matchId", validation.Required)))
	if err != nil {
		u.ShowResponse("Failure", 400, err, w)
		return
	}

	if claims.UserID != matchData.U_ID {
		u.ShowResponse("Failure", 400, "This user did not created the match", w)
		return
	}

	matchData.Status = "Completed"

	//find the scorecard relatedd to that match
	var scorecard models.MatchRecord
	err = db.DB.Where("m_id=?", mp["matchId"].(string)).First(&scorecard).Error
	if err != nil {
		u.ShowResponse("Failure", 400, err, w)
		return
	}

	//Fetch all the score and update the result
	var records []models.ScoreCard
	err = db.DB.Where("s_id", scorecard.S_ID).Find(&records).Error
	if err != nil {
		u.ShowResponse("Failure", 400, err, w)
		return
	}

	for _, record := range records {
		//fmt.Println("Player id is:", record.P_ID)
		var pCareer models.Career
		err = db.DB.Where("p_id=?", record.P_ID).First(&pCareer).Error
		if err != nil {
			u.ShowResponse("Failure", 400, err, w)
			return
		}

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
		err = db.DB.Where("p_id=?", record.P_ID).Updates(&pCareer).Error
		if err != nil {
			u.ShowResponse("Failure", 400, err, w)
			return
		}
	}

	var teamsRuns []models.Inning
	err = db.DB.Where("m_id=?", mp["matchId"].(string)).Find(&teamsRuns).Error
	if err != nil {
		u.ShowResponse("Failure", 400, err, w)
		return
	}

	fmt.Println("jdfj:", teamsRuns)
	if teamsRuns[0].TScore > teamsRuns[1].TScore {
		matchData.Text = teamsRuns[0].T_ID + " Won the match"
	} else {
		matchData.Text = teamsRuns[1].T_ID + " Won the match"
	}

	err = db.DB.Where("m_id=?", mp["matchId"].(string)).Updates(&matchData).Error
	if err != nil {
		u.ShowResponse("Failure", 400, err, w)
		return
	}

	u.ShowResponse("Success", http.StatusOK, matchData, w)

}

// @Description Shows a particular match
// @Accept json
// @Produces json
// @Success 200 {object} models.Response
// @Param matchId body object true "Match Id"
// @Tags Match
// @Router /showMatchById [post]
func ShowMatchById(w http.ResponseWriter, r *http.Request) {
	var mp = make(map[string]string)
	var match models.Match
	err := json.NewDecoder(r.Body).Decode(&mp)
	if err != nil {
		u.ShowResponse("Failure", 400, err, w)
		return
	}
	matchId := mp["matchId"]
	if matchId == "" {
		u.ShowResponse("Failure", 400, "Plaease provide match id", w)
		return
	}

	err = db.DB.Where("m_id=?", matchId).First(&match).Error
	if err != nil {
		u.ShowResponse("Failure", 400, err, w)
		return
	}

	u.ShowResponse("Success", 200, match, w)
}
