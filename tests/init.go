package tests

import (
	"log"

	"github.com/joho/godotenv"
)

var database *TDB

func init() {
	if err := godotenv.Load("test.env"); err != nil {
		log.Print("No .env file found")
	}
	database = NewFromEnv()
}
