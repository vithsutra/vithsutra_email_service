package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/vithsutra/vithsutra_email_sending_service/consumer"
)

func init() {
	serverMode := os.Getenv("SERVER_MODE")

	if serverMode == "dev" {
		if err := godotenv.Load(); err != nil {
			log.Fatalln(".env file was missing, failed to load")
		}
		return
	}

	if serverMode == "prod" {
		return
	}

	log.Fatalln("please set SERVER_MODE to dev or prod")
}

func main() {
	consumer.Start()
}
