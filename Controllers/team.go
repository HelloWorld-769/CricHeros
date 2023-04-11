package controllers

import (
	db "cricHeros/Database"
	models "cricHeros/Models"
	u "cricHeros/Utils"
	"encoding/json"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// @Description Creates a team
// @Accept json
// @Produces json
// @Success 200 {object} models.Response
// @Param id query string true "ID of the user"
// @Param TeamDetails body models.Team true "Details of the team"
// @Tags Team
// @Router /createTeam [post]
func CreateTeamHandler(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")
	u.SetHeader(w)
	u.EnableCors(&w)
	var team models.Team
	err := json.NewDecoder(r.Body).Decode(&team)

	if err != nil {
		u.ShowResponse("Failure", 400, err, w)
		return
	}

	validationErr := u.CheckValidation(team)
	if validationErr != nil {
		u.ShowResponse("Failure", 400, err, w)
		return
	}

	team.U_ID = id
	err = db.DB.Create(&team).Error
	if err != nil {
		u.ShowResponse("Failure", 400, err, w)
		return
	}

	u.ShowResponse("Success", 200, team, w)
}

// @Description Add player to team
// @Accept json
// @Success 200 {object} models.Response
// @Param id query string true "ID of the team"
// @Param player body []string true "Array of players"
// @Tags Team
// @Router /addPlayertoTeam [post]
func AddPlayertoTeamHandler(w http.ResponseWriter, r *http.Request) {
	u.SetHeader(w)
	u.EnableCors(&w)
	var mp = make(map[string][]string)
	var team models.Team
	err := json.NewDecoder(r.Body).Decode(&mp)
	if err != nil {
		u.ShowResponse("Failure", 400, err, w)
		return
	}

	teamID := mp["teamId"][0]
	if teamID == "" {
		u.ShowResponse("Failure", 400, "Please provide team id", w)
		return
	}

	err = db.DB.Where("t_id=?", teamID).First(&team).Error
	if err != nil {
		u.ShowResponse("Failure", 400, err, w)
		return
	}

	for _, p := range mp["players"] {

		var exists bool
		query := "SELECT EXISTS(SELECT * FROM teams where p_id=?)"

		db.DB.Raw(query, p).Scan(&exists)
		if exists {
			u.ShowResponse("Failure", 400, "Player already present in that team", w)
			return
		}
		var player models.Player
		team.P_ID = p
		db.DB.Where("p_id=?", teamID).Updates(&player)
		teamList := models.TeamList{
			P_ID: p,
			T_ID: teamID,
		}
		err := db.DB.Create(&teamList).Error
		if err != nil {
			u.ShowResponse("Failure", 500, "Internal Server Error", w)
			return
		}
		err = db.DB.Create(&team).Error
		if err != nil {
			u.ShowResponse("Failure", 500, "Internal Server Error", w)
			return
		}
	}
	err = db.DB.Exec("DELETE FROM teams WHERE p_id='' and t_id=?", teamID).Error
	if err != nil {
		u.ShowResponse("Failure", 400, err, w)
		return
	}

	u.ShowResponse("Success", 200, "Players added to the team", w)

}

// @Description Shows the list of teams
// @Accept json
// @Produces json
// @Success 200 {object} models.Response
// @Param id query string true "ID of the User"
// @Tags Team
// @Router /showTeams [get]
func ShowTeamsHandler(w http.ResponseWriter, r *http.Request) {
	u.EnableCors(&w)
	u.SetHeader(w)
	var mp = make(map[string]string)
	var teams []models.Team
	err := json.NewDecoder(r.Body).Decode(&mp)
	if err != nil {
		u.ShowResponse("Failure", 400, err, w)
		return
	}
	err = validation.Validate(mp,
		validation.Map(
			validation.Key("userId", validation.Required),
		),
	)
	if err != nil {
		u.ShowResponse("Failure", 400, err, w)
		return
	}
	id := mp["userId"]

	query := "SELECT DISTINCT t_id,t_name,t_captain,t_type FROM teams where u_id=?"
	err = db.DB.Raw(query, id).Scan(&teams).Error
	if err != nil {
		u.ShowResponse("Failure", 400, err, w)
		return
	}

	u.ShowResponse("Success", 200, teams, w)

}

// @Description Shows the list of teams
// @Accept json
// @Produces json
// @Success 200 {object} models.Response
// @Param  team_id body string  true "ID of the team".
// @Tags Team
// @Router /showTeamByID [post]
func ShowTeamByIDHandler(w http.ResponseWriter, r *http.Request) {
	u.SetHeader(w)
	u.EnableCors(&w)
	var mp = make(map[string]string)
	var team models.Team
	var player []string
	err := json.NewDecoder(r.Body).Decode(&mp)
	if err != nil {
		u.ShowResponse("Failure", 400, err, w)
		return
	}

	t_id := mp["teamId"]
	u_id := mp["userId"]

	//var players []string
	query := "SELECT t_id,t_name,t_captain,t_type FROM teams WHERE t_id =?"
	err = db.DB.Raw(query, t_id).Scan(&team).Error
	if err != nil {
		u.ShowResponse("Failure", 400, err, w)
		return
	}

	query = "SELECT p.p_name from players as p JOIN teams as t ON p.p_id=t.p_id WHERE t.t_id=? AND t.u_id=?"
	err = db.DB.Raw(query, t_id, u_id).Scan(&player).Error
	if err != nil {
		u.ShowResponse("Failure", 400, err, w)
		return
	}

	u.ShowResponse("Success", 200, &team, w)
	u.ShowResponse("Success", 200, &player, w)

}

// @Description Delete the team
// @Accept json
// @Produces json
// @Success 200 {object} models.Response
// @Param id query string true "ID of the team".
// @Param user_id body object true "ID of the user"
// @Tags Team
// @Router /deleteTeamByID [delete]
func DeleteTeamHandler(w http.ResponseWriter, r *http.Request) {
	u.SetHeader(w)
	u.EnableCors(&w)
	id := r.URL.Query().Get("id")
	if id == "" {
		u.ShowResponse("Failure", 400, "Id not provided", w)
		return
	}

	query := "DELETE FROM teams WHERE t_id=?;"
	err := db.DB.Raw(query, id).Error
	if err != nil {
		u.ShowResponse("Failure", 400, err, w)
		return
	}

	u.ShowResponse("Success", 200, "Team deleted successfully", w)
}
