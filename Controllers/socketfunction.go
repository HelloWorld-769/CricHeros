package controllers

import (
	db "cricHeros/Database"
	models "cricHeros/Models"

	socketio "github.com/googollee/go-socket.io"
)

func GetData(s socketio.Conn, mp map[string]string) {
	var matchDetail models.MatchRecord
	var scoreCardRecord []models.ScoreCard
	db.DB.Where("m_id=?", mp["matchId"]).First(&matchDetail)

	db.DB.Where("s_id=?", matchDetail.S_ID).Find(&scoreCardRecord)

	s.Emit("scorecard", scoreCardRecord)

}
