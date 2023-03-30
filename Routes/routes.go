package routes

import (
	c "cricHeros/Controllers"
	db "cricHeros/Database"
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

	//Player routes
	mux.HandleFunc("/createPlayer", c.AddPlayerHandler)
	mux.HandleFunc("/showPlayer", c.ShowPlayerHandler)
	mux.HandleFunc("/showPlayerID", c.ShowPlayerByIDHandler)

	//Career routes
	mux.HandleFunc("/addCareer", c.AddCareerHandler)

	//team routes
	mux.HandleFunc("/createTeam", c.CreateTeamHandler)
	mux.HandleFunc("/addPlayertoTeam", c.AddPlayertoTeamHandler)
	mux.HandleFunc("/showTeams", c.ShowTeamsHandler)
	mux.HandleFunc("/showTeamByID", c.ShowTeamByIDHandler)

	//Authentication Handler
	mux.HandleFunc("/register", c.RegisterHandler)
	mux.HandleFunc("/login", c.LoginHandler)
	mux.HandleFunc("/forgotPassword", c.ForgotPasswordHandler)
	mux.HandleFunc("/resetPassword", c.ResetPasswordHandler)

	//Match routes
	mux.HandleFunc("/createMatch", c.CreateMatchHandler)
	mux.HandleFunc("/showMatch", c.ShowMatchHandler)
	mux.HandleFunc("/endMatch", c.EndMatchHandler)

	//score card routes
	mux.HandleFunc("/addToScoreCard", c.ScorecardRecordHandler)

	log.Fatal(http.ListenAndServe(os.Getenv("PORT"), mux))
}
