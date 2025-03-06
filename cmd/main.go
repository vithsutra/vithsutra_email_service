//boiler plate code

package main

import (
	"email-service/config"
	"email-service/internal/kafka"
	"log"
)

func main() {
	// Load Config
	err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Start Kafka Consumer
	kafka.StartConsumer()
}
