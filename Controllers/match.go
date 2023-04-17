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
	userId := r.Context().Value("userId")
	var match models.Match
	var count int64
	var exists bool
	var exists2 bool
	err := json.NewDecoder(r.Body).Decode(&match)
	if err != nil {
		u.ShowResponse("Failure", 400, err, w)
		return
	}

	//check if the entered team has no player then show error

	//check if the entered teams ids exists or not
	query := "SELECT EXISTS(SELECT * FROM teams WHERE t_id=?)"
	db.DB.Raw(query, match.T1_ID).Scan(&exists)
	db.DB.Raw(query, match.T2_ID).Scan(&exists2)
	if !exists || !exists2 {
		u.ShowResponse("Failure", 400, "Team with given team id do not exists please register the team first", w)
		return
	}

	query = "SELECT EXISTS(SELECT * FROM teams WHERE t_id=? AND p_id='')"
	db.DB.Raw(query, match.T1_ID).Scan(&exists)
	db.DB.Raw(query, match.T2_ID).Scan(&exists2)
	if exists || exists2 {
		u.ShowResponse("Failure", 400, "The selected team has no players", w)
		return
	}

	if match.T1_ID == match.T2_ID {
		u.ShowResponse("Failure", 400, "Team 1 and Team 2 cannot be same", w)
		return
	}
	//check if the match is created with same teams but the previous match is not completed yet
	query = `SELECT EXISTS(SELECT status FROM matches where t1_id=? AND t2_id=? AND status='active')`
	db.DB.Raw(query, match.T1_ID, match.T2_ID).Scan(&exists)
	if exists {
		u.ShowResponse("Failure", 400, "Match with the given team is not completed yet.Please complete that match first.", w)
		return
	}

	//check if one player is common in both teams
	query = "SELECT * FROM teams where t_id=? and p_id IN (SELECT p_id FROM teams WHERE t_id=? )"
	db.DB.Raw(query, match.T1_ID, match.T2_ID).Count(&count)
	if count != 0 {
		u.ShowResponse("Failure", 400, "Common record found please edit your teams", w)
		return
	}

	match.U_ID = userId.(string)
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
// @Param matchDetils body string true "Id of the match to end it" SchemaExample({\n "matchId":"string",\n "teamId":"string"\n})
// @Tags Match
// @Router /endMatch [post]
func EndMatchHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("dsfbjsabdjfbjd")
	u.EnableCors(&w)
	u.SetHeader(w)
	var mp = make(map[string]interface{})
	var matchData models.Match
	var scorecard models.MatchRecord
	var records []models.ScoreCard
	var teamsRuns []models.Inning
	var exists bool
	userId := r.Context().Value("userId")

	err := json.NewDecoder(r.Body).Decode(&mp)
	if err != nil {
		u.ShowResponse("Failure1", 400, err, w)
		return
	}

	//validaton check
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

	//Check if the match exists or not
	query := "SELECT EXISTS(SELECT * FROM matches where m_id=?)"
	db.DB.Raw(query, mp["matchId"].(string)).Scan(&exists)
	if !exists {
		u.ShowResponse("Failure", 400, "Match do not exists", w)
		return
	}
	EndInningHandler2(mp, w)

	//get match data to update its status to completed
	err = db.DB.Where("m_id", mp["matchId"].(string)).First(&matchData).Error
	if err != nil {
		u.ShowResponse("Failure", 400, err.Error(), w)
		return
	}

	//check if the match is ended by that user only who created it
	if matchData.U_ID != userId {
		u.ShowResponse("Failure", 400, "This user did not created the match", w)
		return
	}
	//check if the match is already completed or not
	if matchData.Status == "Completed" {
		u.ShowResponse("Failure", 400, "Match has already been completed", w)
		return
	}

	matchData.Status = "completed"

	//find the scorecard relatedd to that match
	err = db.DB.Where("m_id=?", mp["matchId"].(string)).First(&scorecard).Error
	if err != nil {
		u.ShowResponse("Failure", 400, err, w)
		return
	}

	//Fetch all the score and update the result

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
			u.ShowResponse("Failure", 400, err.Error(), w)
			return
		}

		pCareer.MPlayed += 1
		if record.PType == "batsmen" {
			pCareer.RunScored += record.RunScored
			pCareer.Fours += record.Fours
			pCareer.Sixes += record.Sixes
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

	err = db.DB.Where("m_id=?", mp["matchId"].(string)).Find(&teamsRuns).Error
	if err != nil {
		u.ShowResponse("Failure", 400, err, w)
		return
	}

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
// @Param matchDetails body string true "Id of the match " SchemaExample({\n "matchId":"string" \n})
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
	err = validation.Validate(mp,
		validation.Map(
			validation.Key("matchId", validation.Required),
		),
	)
	if err != nil {
		u.ShowResponse("Failure", 400, err, w)
		return
	}
	matchId := mp["matchId"]

	err = db.DB.Where("m_id=?", matchId).First(&match).Error
	if err != nil {
		u.ShowResponse("Failure", 400, err, w)
		return
	}

	u.ShowResponse("Success", 200, match, w)
}

// @Description Deletes the match
// @Accept json
// @Produces json
// @Success 200 {object} models.Response
// @Param matchId body string true "Match Id" ScehmaExample({\n"matchId":"string"\n})
// @Tags Match
// @Router /deleteMatch [delete]
func DeleteMatchHandler(w http.ResponseWriter, r *http.Request) {
	var mp = make(map[string]string)
	err := json.NewDecoder(r.Body).Decode(&mp)
	if err != nil {
		u.ShowResponse("Failure", 400, "Unable to decode the json objec", w)
		return
	}

	err = validation.Validate(mp,
		validation.Map(
			validation.Key("matchId", validation.Required),
		),
	)
	if err != nil {
		u.ShowResponse("Failure", 400, err, w)
		return
	}
	err = db.DB.Where("m_id=?", mp["matchId"]).Delete(&models.Match{}).Error
	if err != nil {
		u.ShowResponse("Failure", 400, err.Error(), w)
		return
	}
	err = db.DB.Where("m_id=?", mp["matchId"]).Delete(&models.MatchRecord{}).Error
	if err != nil {
		u.ShowResponse("Failure", 400, err.Error(), w)
		return
	}

	u.ShowResponse("Success", 200, "Match deleted successfully", w)
}
