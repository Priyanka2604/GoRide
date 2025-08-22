package main

import (
	"context"
	"driver_svc/internal/handlers"
	"driver_svc/internal/mq"
	"driver_svc/internal/repo"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	mongoURI := getEnv("MONGO_URI", "mongodb://localhost:27017")
	kafkaBroker := getEnv("KAFKA_BROKER", "localhost:9092")

	// Mongo connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("❌ Failed to connect to MongoDB:", err)
	}
	db := client.Database("goride")

	repo := repo.NewDriverRepo(db)
	_ = repo.SeedDrivers(ctx)

	// MQ producer/consumer
	producer := mq.NewProducer(kafkaBroker, "booking.accepted")
	consumer := mq.NewConsumer(kafkaBroker, "booking.created", "driver-svc-group", repo)

	handler := handlers.NewDriverHandler(repo, producer)

	// Router
	r := chi.NewRouter()
	r.Get("/drivers", handler.GetDrivers)
	r.Get("/jobs", handler.GetJobs)
	r.Post("/jobs/{booking_id}/accept", handler.AcceptJob)

	// Start consumer
	go consumer.Start(context.Background())

	log.Println("🚀 driver_svc running on :8081")
	if err := http.ListenAndServe(":8081", r); err != nil {
		log.Fatal(err)
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
