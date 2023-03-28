package controllers

import (
	db "cricHeros/Database"
	models "cricHeros/Models"
	"encoding/json"
	"net/http"
)

func CreateTeamHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	var team models.Team
	json.NewDecoder(r.Body).Decode(&team)
	db.DB.Create(&team)
	json.NewEncoder(w).Encode(&team)
}
func AddPlayertoTeamHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var mp = make(map[string][]string)

	json.NewDecoder(r.Body).Decode(&mp)
	id := r.URL.Query().Get("id")
	var team models.Team
	db.DB.Where("t_id=?", id).First(&team)
	for _, p := range mp["players"] {
		team.P_ID = p
		db.DB.Create(&team)
	}
	db.DB.Exec("DELETE FROM teams WHERE p_id=''")

}
func ShowTeamsHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	var teams []models.Team
	query := "SELECT DISTINCT t_id,t_name,t_captain,t_type FROM teams"
	db.DB.Raw(query).Scan(&teams)

	json.NewEncoder(w).Encode(teams)
}

func ShowTeamByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	var mp = make(map[string]string)
	json.NewDecoder(r.Body).Decode(&mp)
	id := mp["id"]

	var team models.Team
	var player []string
	//var players []string
	query := "SELECT t_id,t_name,t_captain,t_type FROM teams WHERE t_id =?"
	db.DB.Raw(query, id).Scan(&team)

	query = "SELECT p.p_name from players as p JOIN teams as t ON p.p_id=t.p_id WHERE t.t_id=?"
	db.DB.Raw(query, id).Scan(&player)

	json.NewEncoder(w).Encode(&team)
	json.NewEncoder(w).Encode(&player)

}
