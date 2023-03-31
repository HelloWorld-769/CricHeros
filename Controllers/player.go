package controllers

import (
	db "cricHeros/Database"
	models "cricHeros/Models"
	u "cricHeros/Utils"
	"encoding/json"
	"fmt"
	"net/http"
)

func AddPlayerHandler(w http.ResponseWriter, r *http.Request) {
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

// function to add players batting career
func ShowPlayerHandler(w http.ResponseWriter, r *http.Request) {

	//teams mapping left
	u.SetHeader(w)
	var players []models.Player
	err := db.DB.Find(&players).Error
	if err != nil {
		fmt.Println("Error in extracting the data from the database", err)
		return
	}

	u.Encode(w, &players)
}

func ShowPlayerByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var PlayerData models.Response
	var mp = make(map[string]string)
	json.NewDecoder(r.Body).Decode(&mp)
	id := mp["id"]

	db.DB.Table("players").Where("p_id=?", id).Scan(&PlayerData.Player)

	db.DB.Table("careers").Where("p_id=?", id).Scan(&PlayerData.Career)

	db.DB.Table("team_lists").Where("p_id=?", id).Scan(&PlayerData.Teams)

	u.Encode(w, &PlayerData)
}
