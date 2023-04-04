package controllers

import (
	db "cricHeros/Database"
	models "cricHeros/Models"
	u "cricHeros/Utils"
	"encoding/json"
	"fmt"
	"net/http"
)

// @Description Creates a new Player
// @Accept json
// @Produce json
// @Tags Player
// @Param player body models.Player true "Create Player"
// @Success 200 {object} models.Player
// @Router /createPlayer [post]
func AddPlayerHandler(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	u.SetHeader(w)
	var player models.Player
	err := json.NewDecoder(r.Body).Decode(&player)
	if err != nil {
		fmt.Println("Error in decoding the body")
		return
	}
	db.DB.Create(&player)
	if err != nil {
		fmt.Println("Error in inserting data: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	u.Encode(w, &player)

}

// @Description Shows the list of all the player
// @Accept json
// @Produce json
// @Success 200 {object} models.Player
// @Failure 400 {string} string "Bad Request"
// @Tags Player
// @Router /showPlayer [get]
func ShowPlayerHandler(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	u.SetHeader(w)
	var players []models.Player
	err := db.DB.Find(&players).Error
	if err != nil {
		fmt.Println("Error in extracting the data from the database", err)
		return
	}

	u.Encode(w, &players)
}

// @Description Shows the list of all the player
// @Accept json
// @Produce json
// @Success 200 {object} models.Response
// @Param id body object true "ID of the player"
// @Tags Player
// @Router /showPlayerID [get]
func ShowPlayerByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	EnableCors(&w)
	var PlayerData models.Response
	var mp = make(map[string]string)
	json.NewDecoder(r.Body).Decode(&mp)
	id := mp["id"]

	db.DB.Table("players").Where("p_id=?", id).Scan(&PlayerData.Player)

	db.DB.Table("careers").Where("p_id=?", id).Scan(&PlayerData.Career)

	db.DB.Table("team_lists").Where("p_id=?", id).Scan(&PlayerData.Teams)

	u.Encode(w, &PlayerData)
}
