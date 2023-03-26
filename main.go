package main

import (
	"fmt"

	r "cricHeros/Routes"

	"github.com/joho/godotenv"
)

func main() {
	envErr := godotenv.Load(".env")
	if envErr != nil {
		fmt.Println("Could not load environment variable")
		return
	}
	r.Routes()

}
