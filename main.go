package main

import (
	"fmt"

	r "cricHeros/Routes"

	"github.com/joho/godotenv"
)

// @title Cric Heros API
// @version 1.0.0
// @description API Documentation for Cric Heros
// @host localhost:8000
func main() {
	envErr := godotenv.Load(".env")
	if envErr != nil {
		fmt.Println("Could not load environment variable")
		return
	}
	r.Routes()

}
