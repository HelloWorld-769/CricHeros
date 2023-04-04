package controllers

import (
	db "cricHeros/Database"
	models "cricHeros/Models"
	u "cricHeros/Utils"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"
)

var tossOptions = []string{"Head", "Tails"}

func getRandomResult() string {
	rand.Seed(time.Now().UnixNano())
	result := tossOptions[rand.Intn(len(tossOptions))]
	return result
}

// @Description Give the random result of coin toss and which team won the toss
// @Accept json
// @Produces json
// @Success 200 {object} models.Toss
// @Param toss body models.Toss true "Toss Details"
// @Tags Toss
// @Router /tossResult [post]
func TossResultHandler(w http.ResponseWriter, r *http.Request) {
	u.SetHeader(w)
	EnableCors(&w)
	var toss models.Toss
	json.NewDecoder(r.Body).Decode(&toss)
	tossRes := getRandomResult()
	var team_id string
	if tossRes == "Head" {
		team_id = toss.T1_ID
	} else {
		team_id = toss.T2_ID
	}
	res := team_id + " Won the toss"
	toss.TossWon = res
	//db.DB.Create(toss)
	u.Encode(w, &toss)
}

// @Description Updates the decison taken by the team after wining the toss
// @Accept json
// @Produces json
// @Success 200 {object} models.Toss
// @Param toss body models.Toss true "Descision Updated"
// @Tags Toss
// @Router /DecisionUpdate [put]
func DecisionUpdateHandler(w http.ResponseWriter, r *http.Request) {
	u.SetHeader(w)
	EnableCors(&w)
	var tossDecision models.Toss
	json.NewDecoder(r.Body).Decode(&tossDecision)
	db.DB.Where("toss_id=?", tossDecision.Toss_ID).Updates(&tossDecision)
	u.Encode(w, &tossDecision)

}
