package routes

import (
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
	mux.HandleFunc("/welcome", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welomce to Cric heros application"))
	})
	log.Fatal(http.ListenAndServe(os.Getenv("PORT"), mux))
}
