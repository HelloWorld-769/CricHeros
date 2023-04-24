package routes

import (
	c "cricHeros/Controllers"
	db "cricHeros/Database"
	constants "cricHeros/Utils"
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
	mux.Use(c.CORSMiddleware)

	//Authentication Handler
	mux.HandleFunc("/adminRegister", c.AdminRegisterHandler).Methods("POST")
	mux.HandleFunc("/userRegister", c.UserRegisterHandler).Methods("POST")
	mux.HandleFunc("/sendOTP", c.SendOtpHandler).Methods("POST")
	mux.HandleFunc("/verifyOTP", c.VerifyOTPHandler).Methods("POST")
	mux.HandleFunc("/logout", c.LogOut).Methods("GET")
	mux.Handle("/updateProfile", c.LoginMiddlerware(http.HandlerFunc(c.UpdateProfile))).Methods("POST")

	//Player routes
	mux.Handle("/createPlayer", c.LoginMiddlerware(http.HandlerFunc(c.AddPlayerHandler))).Methods("POST")
	mux.Handle("/showPlayer", c.LoginMiddlerware(http.HandlerFunc(c.ShowPlayerHandler))).Methods("GET")
	mux.Handle("/showPlayerID", c.LoginMiddlerware(http.HandlerFunc(c.ShowPlayerByIDHandler))).Methods("POST")
	mux.Handle("/retirePlayer", c.LoginMiddlerware(http.HandlerFunc(c.DeletePlayerHandler))).Methods("DELETE")

	//Career routes
	mux.Handle("/addCareer", c.LoginMiddlerware(http.HandlerFunc(c.AddCareerHandler))).Methods("POST")

	//team routes
	mux.Handle("/createTeam", c.LoginMiddlerware(http.HandlerFunc(c.CreateTeamHandler))).Methods("POST")
	mux.Handle("/addPlayertoTeam", c.LoginMiddlerware(http.HandlerFunc(c.AddPlayertoTeamHandler))).Methods("POST")
	mux.Handle("/showTeams", c.LoginMiddlerware(http.HandlerFunc(c.ShowTeamsHandler))).Methods("GET")
	mux.Handle("/showTeamByID", c.LoginMiddlerware(http.HandlerFunc(c.ShowTeamByIDHandler))).Methods("POST")
	mux.Handle("/deleteTeamByID", c.LoginMiddlerware(http.HandlerFunc(c.DeleteTeamHandler))).Methods("DELETE")

	//Match routes
	mux.Handle("/createMatch", c.AdminMiddlerware(c.LoginMiddlerware(http.HandlerFunc(c.CreateMatchHandler)))).Methods("POST")
	mux.Handle("/showMatch", c.LoginMiddlerware(http.HandlerFunc(c.ShowMatchHandler))).Methods("GET")
	mux.Handle("/endMatch", c.AdminMiddlerware(c.LoginMiddlerware(http.HandlerFunc(c.EndMatchHandler)))).Methods("POST")
	mux.Handle("/showMatchById", c.LoginMiddlerware(http.HandlerFunc(c.ShowMatchById))).Methods("POST")
	mux.Handle("/deleteMatch", c.AdminMiddlerware(c.LoginMiddlerware(http.HandlerFunc(c.DeleteMatchHandler)))).Methods("DELETE")

	//score card routes
	mux.Handle("/addToScoreCard", c.AdminMiddlerware(c.LoginMiddlerware(http.HandlerFunc(c.ScorecardRecordHandler)))).Methods("POST")

	//Innings routes
	mux.Handle("/endInning", c.AdminMiddlerware(c.LoginMiddlerware(http.HandlerFunc(c.EndInningHandler)))).Methods("POST")

	//Toss Routes
	mux.Handle("/tossResult", c.AdminMiddlerware(c.LoginMiddlerware(http.HandlerFunc(c.TossResultHandler)))).Methods("POST")
	mux.Handle("/decisionUpdate", c.AdminMiddlerware(c.LoginMiddlerware(http.HandlerFunc(c.DecisionUpdateHandler)))).Methods("POST")

	//Ball Handler
	mux.Handle("/ballUpdate", c.AdminMiddlerware(c.LoginMiddlerware(http.HandlerFunc(c.UpdateBallRecord)))).Methods("POST")

	//Socket Server

	c.SocketHandler(constants.SocketServer)
	go constants.SocketServer.Serve()
	mux.Handle("/socket.io/", constants.SocketServer)

	//API doucmentation route
	mux.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler).Methods("GET")

	//Listening to the server
	log.Fatal(http.ListenAndServe(os.Getenv("PORT"), mux))
}
