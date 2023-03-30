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

func CreateMatchHandler(w http.ResponseWriter, r *http.Request) {
	var match models.Match
	json.NewDecoder(r.Body).Decode(&match)
	query := "SELECT * FROM teams where t_id=? and p_id IN (SELECT p_id FROM teams WHERE t_id=? )"
	err := db.DB.Exec(query, match.T1_ID, match.T2_ID).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		u.ShowErr("Common record found Please Edit your teams", 400, w)
		return
	}
	matchRecord := models.MatchRecord{
		M_ID: match.M_ID,
		S_ID: match.S_ID,
	}
	db.DB.Create(&matchRecord)
	db.DB.Create(&match)
	u.Encode(w, &match)
}
func ShowMatchHandler(w http.ResponseWriter, r *http.Request) {
	var matches []models.Match
	db.DB.Find(&matches)
	u.Encode(w, &matches)
}

func EndMatchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("End match handler called..."))

	var mp = make(map[string]string)
	json.NewDecoder(r.Body).Decode(&mp)
	var records []models.ScoreCard
	db.DB.Where("s_id", mp["scorecard_id"]).Find(&records)

	fmt.Println("records in scorecard is:", records)

}
