package controllers

import (
	"fmt"

	socketio "github.com/googollee/go-socket.io"
)

func SocketHandler(server *socketio.Server) {
	server.OnConnect("/", func(s socketio.Conn) error {
		fmt.Println("connected:", s.ID())
		return nil
	})
	server.OnEvent("/", "scorecard", GetData)

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})

}
