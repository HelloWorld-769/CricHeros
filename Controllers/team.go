package controllers

import (
	db "cricHeros/Database"
	models "cricHeros/Models"
	u "cricHeros/Utils"
	"encoding/json"
	"fmt"
	"net/http"
)

// @Description Creates a team
// @Accept json
// @Produces json
// @Success 200 {object} models.Team
// @Param id query string true "ID of the user"
// @Param TeamDetails body models.Team true "Details of the team"
// @Tags Team
// @Router /createTeam [post]
func CreateTeamHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	id := r.URL.Query().Get("id")
	u.SetHeader(w)
	var team models.Team
	json.NewDecoder(r.Body).Decode(&team)
	team.U_ID = id
	db.DB.Create(&team)
	json.NewEncoder(w).Encode(&team)
}

// @Description Add player to team
// @Accept json
// @Success 200
// @Param id query string true "ID of the team"
// @Param player body []string true "Array of players"
// @Tags Team
// @Router /addPlayertoTeam [post]
func AddPlayertoTeamHandler(w http.ResponseWriter, r *http.Request) {
	u.SetHeader(w)
	var mp = make(map[string][]string)

	json.NewDecoder(r.Body).Decode(&mp)
	id := r.URL.Query().Get("id")
	var team models.Team
	db.DB.Where("t_id=?", id).First(&team)
	for _, p := range mp["players"] {
		var player models.Player
		team.P_ID = p
		db.DB.Where("p_id=?", id).Updates(&player)
		teamList := models.TeamList{
			P_ID: p,
			T_ID: id,
		}
		db.DB.Create(&teamList)
		db.DB.Create(&team)
	}
	db.DB.Exec("DELETE FROM teams WHERE p_id='' and t_id=?", id)

}

// @Description Shows the list of teams
// @Accept json
// @Produces json
// @Success 200 {object} models.Team
// @Param id query string true "ID of the User"
// @Tags Team
// @Router /showTeams [get]
func ShowTeamsHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var mp = make(map[string]string)
	json.NewDecoder(r.Body).Decode(&mp)
	id := mp["id"]
	u.SetHeader(w)
	var teams []models.Team
	query := "SELECT DISTINCT t_id,t_name,t_captain,t_type FROM teams where u_id=?"
	db.DB.Raw(query, id).Scan(&teams)

	json.NewEncoder(w).Encode(teams)
}

// @Description Shows the list of teams
// @Accept json
// @Produces json
// @Success 200 {object} models.Team
// @Success 200 {object} models.Player
// @Param  team_id body string  true "ID of the team".
// @Tags Team
// @Router /showTeamByID [post]
func ShowTeamByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	u.SetHeader(w)
	var mp = make(map[string]string)
	json.NewDecoder(r.Body).Decode(&mp)
	t_id := mp["team_id"]
	u_id := mp["user_id"]
	var team models.Team
	var player []string
	//var players []string
	query := "SELECT t_id,t_name,t_captain,t_type FROM teams WHERE t_id =?"
	db.DB.Raw(query, t_id).Scan(&team)

	query = "SELECT p.p_name from players as p JOIN teams as t ON p.p_id=t.p_id WHERE t.t_id=? AND t.u_id=?"
	db.DB.Raw(query, t_id, u_id).Scan(&player)

	json.NewEncoder(w).Encode(&team)
	json.NewEncoder(w).Encode(&player)

}

// @Description Delete the team
// @Accept json
// @Produces json
// @Success 200 {string} Team deleted successfull
// @Failure 500 {string}  Unable to delete the tea
// @Param id query string true "ID of the team".
// @Param user_id body object true "ID of the user"
// @Tags Team
// @Router /deleteTeamByID [delete]
func DeleteTeamHandler(w http.ResponseWriter, r *http.Request) {
	u.SetHeader(w)
	EnableCors(&w)
	id := r.URL.Query().Get("id")
	query := "DELETE FROM teams WHERE t_id=?;"
	err := db.DB.Raw(query, id).Error
	if err != nil {
		u.ShowErr("unable to delete team", 500, w)
		return
	}
	fmt.Fprintf(w, "Team deleted successfully")
}
