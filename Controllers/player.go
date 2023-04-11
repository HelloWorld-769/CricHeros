package controllers

import (
	db "cricHeros/Database"
	models "cricHeros/Models"
	u "cricHeros/Utils"
	"encoding/json"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// @Description Creates a new Player
// @Accept json
// @Produce json
// @Tags Player
// @Param player body models.Player true "Create Player"
// @Success 200 {object} models.Response
// @Router /createPlayer [post]
func AddPlayerHandler(w http.ResponseWriter, r *http.Request) {
	u.EnableCors(&w)
	u.SetHeader(w)

	var exists bool
	var player models.Player
	var playerCareer models.Career
	err := json.NewDecoder(r.Body).Decode(&player)

	if err != nil {
		u.ShowResponse("Failure", 400, err, w)
		return
	}

	validationErr := u.CheckValidation(player)
	if validationErr != nil {
		u.ShowResponse("Failure", 400, validationErr, w)
		return
	}

	query := "SELECT EXISTS(SELECT * FROM players WHERE phone_no=?)"
	db.DB.Raw(query, player.PhoneNo).Scan(&exists)
	if exists {
		u.ShowResponse("Failure", 400, "Player already exists", w)
		return
	}
	err = db.DB.Create(&player).Error
	if err != nil {
		u.ShowResponse("Failure", 400, err.Error(), w)
		return
	}

	playerCareer.P_ID = player.P_ID
	err = db.DB.Create(&playerCareer).Error
	if err != nil {
		u.ShowResponse("Failure", 400, err.Error(), w)
		return
	}

	u.ShowResponse("Success", http.StatusOK, &player, w)
}

// @Description Shows the list of all the player
// @Accept json
// @Produce json
// @Success 200 {object} models.Response
// @Failure 400 {string} string "Bad Request"
// @Tags Player
// @Router /showPlayer [get]
func ShowPlayerHandler(w http.ResponseWriter, r *http.Request) {
	u.EnableCors(&w)
	u.SetHeader(w)
	var players []models.Player
	err := db.DB.Find(&players).Error
	if err != nil {
		u.ShowResponse("Failure", 400, err, w)
		return
	}

	u.ShowResponse("Success", http.StatusOK, &players, w)
}

// @Description Shows the list of all the player
// @Accept json
// @Produce json
// @Success 200 {object} models.Response
// @Param id body object true "ID of the player"
// @Tags Player
// @Router /showPlayerID [get]
func ShowPlayerByIDHandler(w http.ResponseWriter, r *http.Request) {
	u.SetHeader(w)
	u.EnableCors(&w)
	var PlayerData models.PlayerData
	var mp = make(map[string]string)
	err := json.NewDecoder(r.Body).Decode(&mp)
	if err != nil {
		u.ShowResponse("Failure", 400, err, w)
		return
	}
	id := mp["playerId"]
	if id == "" {
		u.ShowResponse("Failure", 400, "Please enter player id", w)
		return
	}

	err = db.DB.Table("players").Where("p_id=?", id).Scan(&PlayerData.Player).Error
	if err != nil {
		u.ShowResponse("Failure", 400, err, w)
		return
	}

	err = db.DB.Table("careers").Where("p_id=?", id).Scan(&PlayerData.Career).Error
	if err != nil {
		u.ShowResponse("Failure", 400, err, w)
		return
	}

	err = db.DB.Table("team_lists").Where("p_id=?", id).Scan(&PlayerData.Teams).Error
	if err != nil {
		u.ShowResponse("Failure", 400, err, w)
		return
	}

	u.ShowResponse("Success", http.StatusOK, &PlayerData, w)
}

// @Description Shows the list of all the player
// @Accept json
// @Produce json
// @Success 200 {object} models.Response
// @Param id body object true "Id of the player"
// @Tags Player
// @Router /retirePlayer [delete]
func DeletePlayerHandler(w http.ResponseWriter, r *http.Request) {
	var mp = make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&mp)
	if err != nil {
		u.ShowResponse("Failure", 400, err, w)
		return
	}

	err = validation.Validate(mp,
		validation.Map(
			validation.Key("playerId", validation.Required),
		),
	)
	if err != nil {
		u.ShowResponse("Failure", 400, err, w)
		return
	}

	err = db.DB.Where("p_id=?", mp["playerId"].(string)).Delete(&models.Player{}).Error
	if err != nil {
		u.ShowResponse("Failure", 400, err.Error(), w)
		return
	}

	u.ShowResponse("Success", http.StatusOK, "Player Deleted sucessfullly", w)

}
