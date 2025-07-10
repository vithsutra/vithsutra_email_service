package consumer

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vithsutra/vithsutra_email_sending_service/email"
	"github.com/vithsutra/vithsutra_email_sending_service/internal/models"
)

func Start() {
	redisAddr := os.Getenv("REDIS_ADDRESS")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	queueName := os.Getenv("REDIS_QUEUE_NAME")

	log.Println("Connecting to Redis at:", redisAddr)
	log.Println("Queue Name:", queueName)

	RedisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       0,
	})

	pong, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	}
	log.Println("Redis PING response:", pong)

	for {
		result, err := RedisClient.BRPop(context.Background(), 10*time.Second, queueName).Result()
		if err == redis.Nil {
			continue
		} else if err != nil {
			log.Println("Redis BRPOP error:", err)
			continue
		}

		if len(result) < 2 {
			log.Println("Invalid BRPOP result:", result)
			continue
		}

		message := result[1]
		log.Println("Raw message received:", message)

		emailRequest := new(models.Email)
		if err := json.Unmarshal([]byte(message), emailRequest); err != nil {
			log.Println("Error decoding JSON:", err)
			continue
		}

		switch emailRequest.EmailType {
		case "otp":
			handleWithRetry(email.SendOtpEmail, emailRequest)
		case "welcome":
			handleWithRetry(email.WelcomeEmail, emailRequest)
		default:
			log.Println("Invalid email type received:", emailRequest.EmailType)
		}
	}
}

func handleWithRetry(sendFunc func(*models.Email) error, emailRequest *models.Email) {
	const maxRetries = 3

	for i := 0; i < maxRetries; i++ {
		if err := sendFunc(emailRequest); err == nil {
			log.Printf("[%s] Email sent to: %s", emailRequest.EmailType, emailRequest.To)
			return
		} else {
			log.Printf("[%s] Retry %d/%d failed for: %s, Error: %v",
				emailRequest.EmailType, i+1, maxRetries, emailRequest.To, err)
			time.Sleep(2 * time.Second)
		}
	}

	log.Printf("[%s] Email failed after %d retries: %s", emailRequest.EmailType, maxRetries, emailRequest.To)
}
