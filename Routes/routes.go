package routes

import (
	c "cricHeros/Controllers"
	db "cricHeros/Database"
	socket "cricHeros/Socket"
	_ "cricHeros/docs"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func Routes() {
	fmt.Println("Listening on port", os.Getenv("PORT"))
	mux := mux.NewRouter()
	err := db.Connect()
	if err != nil {
		panic(err)
	}

	socketServer := socket.SocketInit()
	defer socketServer.Close()

	//Player routes
	mux.HandleFunc("/createPlayer", c.AddPlayerHandler).Methods("POST")
	mux.HandleFunc("/showPlayer", c.ShowPlayerHandler).Methods("GET")
	mux.HandleFunc("/showPlayerID", c.ShowPlayerByIDHandler).Methods("POST")
	mux.HandleFunc("/retirePlayer", c.DeletePlayerHandler).Methods("DELETE")

	//Career routes
	mux.HandleFunc("/addCareer", c.AddCareerHandler).Methods("POST")

	//team routes
	mux.HandleFunc("/createTeam", c.CreateTeamHandler).Methods("POST")
	mux.HandleFunc("/addPlayertoTeam", c.AddPlayertoTeamHandler).Methods("POST")
	mux.HandleFunc("/showTeams", c.ShowTeamsHandler).Methods("GET")
	mux.HandleFunc("/showTeamByID", c.ShowTeamByIDHandler).Methods("POST")
	mux.HandleFunc("/deleteTeamByID", c.DeleteTeamHandler).Methods("DELETE")

	//Authentication Handler
	mux.HandleFunc("/adminRegister", c.AdminRegisterHandler).Methods("POST")
	mux.HandleFunc("/userRegister", c.UserRegisterHandler).Methods("POST")
	mux.HandleFunc("/login", c.LoginHandler).Methods("POST")
	mux.HandleFunc("/forgotPassword", c.ForgotPasswordHandler).Methods("POST")
	mux.HandleFunc("/resetPassword", c.ResetPasswordHandler).Methods("POST")
	mux.HandleFunc("/updatePassword", c.UpdatePasswordHandler).Methods("PUT")

	//Match routes
	mux.Handle("/createMatch", c.AdminMiddlerware(http.HandlerFunc(c.CreateMatchHandler))).Methods("POST")
	mux.HandleFunc("/showMatch", c.ShowMatchHandler).Methods("GET")
	mux.Handle("/endMatch", c.AdminMiddlerware(http.HandlerFunc(c.EndMatchHandler))).Methods("POST")

	//score card routes
	mux.Handle("/addToScoreCard", c.AdminMiddlerware(http.HandlerFunc(c.ScorecardRecordHandler))).Methods("POST")

	//Innings routes
	mux.Handle("/endInning", c.AdminMiddlerware(http.HandlerFunc(c.EndInningHandler))).Methods("POST")

	//Toss Routes
	mux.Handle("/tossResult", c.AdminMiddlerware(http.HandlerFunc(c.TossResultHandler))).Methods("POST")
	mux.Handle("/decisionUpdate", c.AdminMiddlerware(http.HandlerFunc(c.DecisionUpdateHandler))).Methods("POST")

	//Ball Handler
	mux.Handle("/ballUpdate", c.AdminMiddlerware(http.HandlerFunc(c.UpdateBallRecord))).Methods("POST")

	mux.HandleFunc("/showMatchById", c.ShowMatchById).Methods("POST")

	//Socket Server
	mux.Handle("/socket.io/", socketServer)

	//API doucmentation route
	mux.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler).Methods("GET")

	//Listening to the server
	log.Fatal(http.ListenAndServe(os.Getenv("PORT"), mux))
}
