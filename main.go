package main

import (
	twilio "cricHeros/Controllers"
	r "cricHeros/Routes"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// @title Cric Heros API
// @version 1.0.0
// @description API Documentation for Cric Heros
// @host localhost:8000
func main() {
	envErr := godotenv.Load(".env")
	twilio.TwilioInit(os.Getenv("TWILIO_AUTH_TOKEN"))
	if envErr != nil {
		fmt.Println("Could not load environment variable")
		return
	}
	r.Routes()

}
