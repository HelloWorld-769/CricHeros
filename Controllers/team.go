package controllers

import (
	db "cricHeros/Database"
	models "cricHeros/Models"
	u "cricHeros/Utils"
	"encoding/json"
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
	u.EnableCors(&w)
	var team models.Team
	err := json.NewDecoder(r.Body).Decode(&team)

	if err != nil {
		u.ShowResponse("Failure", 400, err.Error(), w)
		return
	}

	validationErr := u.CheckValidation(team)
	if validationErr != nil {
		u.ShowResponse("Failure", 400, validationErr.Error(), w)
		return
	}

	team.U_ID = id
	err = db.DB.Create(&team).Error
	if err != nil {
		u.ShowResponse("Failure", 400, err.Error(), w)
		return
	}
	u.ShowResponse("Success", 200, team, w)
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
	u.EnableCors(&w)
	var mp = make(map[string][]string)

	err := json.NewDecoder(r.Body).Decode(&mp)
	if err != nil {
		u.ShowResponse("Failure", 400, err.Error(), w)
		return
	}
	id := r.URL.Query().Get("id")
	var team models.Team
	err = db.DB.Where("t_id=?", id).First(&team).Error
	if err != nil {
		u.ShowResponse("Failure", 400, err.Error(), w)
		return
	}
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
	u.ShowResponse("Success", 200, "Players added to the team", w)

}

// @Description Shows the list of teams
// @Accept json
// @Produces json
// @Success 200 {object} models.Team
// @Param id query string true "ID of the User"
// @Tags Team
// @Router /showTeams [get]
func ShowTeamsHandler(w http.ResponseWriter, r *http.Request) {

	u.EnableCors(&w)
	u.SetHeader(w)
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var mp = make(map[string]string)
	err := json.NewDecoder(r.Body).Decode(&mp)
	if err != nil {
		u.ShowResponse("Failure", 400, err.Error(), w)
		return
	}
	id := mp["id"]
	var teams []models.Team
	query := "SELECT DISTINCT t_id,t_name,t_captain,t_type FROM teams where u_id=?"
	db.DB.Raw(query, id).Scan(&teams)

	u.ShowResponse("Success", 200, teams, w)

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
	u.EnableCors(&w)
	var mp = make(map[string]string)
	err := json.NewDecoder(r.Body).Decode(&mp)
	if err != nil {
		u.ShowResponse("Failure", 400, err.Error(), w)
		return
	}
	t_id := mp["teamId"]
	u_id := mp["userId"]
	var team models.Team
	var player []string
	//var players []string
	query := "SELECT t_id,t_name,t_captain,t_type FROM teams WHERE t_id =?"
	db.DB.Raw(query, t_id).Scan(&team)

	query = "SELECT p.p_name from players as p JOIN teams as t ON p.p_id=t.p_id WHERE t.t_id=? AND t.u_id=?"
	db.DB.Raw(query, t_id, u_id).Scan(&player)

	u.ShowResponse("Success", 200, &team, w)
	u.ShowResponse("Success", 200, &player, w)

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
	u.EnableCors(&w)
	id := r.URL.Query().Get("id")
	query := "DELETE FROM teams WHERE t_id=?;"
	err := db.DB.Raw(query, id).Error
	if err != nil {
		u.ShowResponse("Success", 500, err, w)
		return
	}

	u.ShowResponse("Success", 200, "Team deleted successfully", w)
}
