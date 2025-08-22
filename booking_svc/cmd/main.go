package main

import (
	"booking_svc/internal/handlers"
	"booking_svc/internal/mq"
	"booking_svc/internal/repo"
	"context"
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
	//topic := "booking.created"

	// Mongo connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("❌ Failed to connect to MongoDB:", err)
	}
	db := client.Database("goride")

	// Repo + Producer
	repo := repo.NewBookingRepo(db)
	producer := mq.NewProducer(kafkaBroker, "booking.created")
	consumer := mq.NewConsumer(kafkaBroker, "booking.accepted", "booking-svc-group", repo)

	handler := handlers.NewBookingHandler(repo, producer)

	// Router
	r := chi.NewRouter()
	r.Post("/bookings", handler.CreateBooking)
	r.Get("/bookings", handler.GetAllBookings)
	r.Get("/bookings/{id}", handler.GetBookingByID)

	// Run consumer in a goroutine
	go consumer.Start(context.Background())

	log.Println("🚀 booking_svc running on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
