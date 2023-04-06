package routes

import (
	c "cricHeros/Controllers"
	db "cricHeros/Database"
	socket "cricHeros/Socket"
	"fmt"
	"log"
	"net/http"
	"os"
)

func Routes() {
	fmt.Println("Listening on port:", os.Getenv("PORT"))
	mux := http.NewServeMux()
	err := db.Connect()
	if err != nil {
		panic(err)
	}

	socketServer := socket.SocketInit()
	defer socketServer.Close()

	//Player routes
	mux.HandleFunc("/createPlayer", c.AddPlayerHandler)
	mux.HandleFunc("/showPlayer", c.ShowPlayerHandler)
	mux.HandleFunc("/showPlayerID", c.ShowPlayerByIDHandler)
	mux.HandleFunc("/retirePlayer", c.DeletePlayerHandler)

	//Career routes
	mux.HandleFunc("/addCareer", c.AddCareerHandler)

	//team routes
	mux.HandleFunc("/createTeam", c.CreateTeamHandler)
	mux.HandleFunc("/addPlayertoTeam", c.AddPlayertoTeamHandler)
	mux.HandleFunc("/showTeams", c.ShowTeamsHandler)
	mux.HandleFunc("/showTeamByID", c.ShowTeamByIDHandler)
	mux.HandleFunc("/deleteTeamByID", c.DeleteTeamHandler)

	//Authentication Handler
	mux.HandleFunc("/adminRegister", c.AdminRegisterHandler)
	mux.HandleFunc("/userRegister", c.UserRegisterHandler)
	mux.HandleFunc("/login", c.LoginHandler)
	mux.HandleFunc("/forgotPassword", c.ForgotPasswordHandler)
	mux.HandleFunc("/resetPassword", c.ResetPasswordHandler)
	mux.HandleFunc("/updatePassword", c.UpdatePasswordHandler)

	//Match routes
	mux.Handle("/createMatch", c.AdminMiddlerware(http.HandlerFunc(c.CreateMatchHandler)))
	mux.HandleFunc("/showMatch", c.ShowMatchHandler)
	mux.Handle("/endMatch", c.AdminMiddlerware(http.HandlerFunc(c.EndMatchHandler)))

	//score card routes
	mux.Handle("/addToScoreCard", c.AdminMiddlerware(http.HandlerFunc(c.ScorecardRecordHandler)))

	//Innings routes
	mux.Handle("/endInning", c.AdminMiddlerware(http.HandlerFunc(c.EndInningHandler)))

	//Toss Routes
	mux.Handle("/tossResult", c.AdminMiddlerware(http.HandlerFunc(c.TossResultHandler)))
	mux.Handle("/DecisionUpdate", c.AdminMiddlerware(http.HandlerFunc(c.DecisionUpdateHandler)))

	//Socket Server
	mux.Handle("/socket.io/", socketServer)

	//Ball Handler
	mux.Handle("/ballUpdate", c.AdminMiddlerware(http.HandlerFunc(c.UpdateBallRecord)))

	//Listening to the server
	log.Fatal(http.ListenAndServe(os.Getenv("PORT"), mux))
}
