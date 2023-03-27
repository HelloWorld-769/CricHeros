package controllers

import (
	db "cricHeros/Database"
	models "cricHeros/Models"
	"encoding/json"
	"fmt"
	"net/http"
)

func AddPlayerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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

	json.NewEncoder(w).Encode(player)

}

// function to add players batting career
func ShowPlayerHandler(w http.ResponseWriter, r *http.Request) {

	//not complete
	w.Header().Set("Content-Type", "application/json")
	var players []models.Player
	err := db.DB.Find(&players).Error
	if err != nil {
		fmt.Println("Error in extracting the data from the database", err)
		return
	}
	json.NewEncoder(w).Encode(&players)
}

func UpdatePlayerHandler(w http.ResponseWriter, r *http.Request) {

}

func UpdatePlayerScoreHandler(w http.ResponseWriter, r *http.Request) {

}

func ShowPlayerByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var PlayerData models.Response
	var mp = make(map[string]string)
	json.NewDecoder(r.Body).Decode(&mp)
	id := mp["id"]

	db.DB.Table("players").Where("p_id=?", id).Scan(&PlayerData.Player)

	db.DB.Table("battings").Where("p_id=?", id).Scan(&PlayerData.Batting)

	db.DB.Table("bowlings").Where("p_id=?", id).Scan(&PlayerData.Bowling)

	json.NewEncoder(w).Encode(&PlayerData)
}
func MakeCaptain(w http.ResponseWriter, r *http.Request) {
	var mp = make(map[string]string)
	json.NewDecoder(r.Body).Decode(&mp)
	id := mp["id"]

	var player models.Player
	db.DB.Where("p_id=?", id).First(&player)
	player.Is_Captain = true

	db.DB.Where("p_id=?", id).Updates(&player)
}
