package mq

import (
	"context"
	"driver_svc/internal/models"
	"driver_svc/internal/repo"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

type BookingCreatedEvent struct {
	BookingID  string          `json:"booking_id"`
	PickupLoc  models.Location `json:"pickuploc"`
	Dropoff    models.Location `json:"dropoff"`
	Price      int             `json:"price"`
	RideStatus string          `json:"ride_status"`
}

type Consumer struct {
	reader *kafka.Reader
	repo   *repo.DriverRepo
}

func NewConsumer(broker, topic, groupID string, r *repo.DriverRepo) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{broker},
			Topic:   topic,
			GroupID: groupID,
		}),
		repo: r,
	}
}

func (c *Consumer) Start(ctx context.Context) {
	log.Println("📥 driver_svc consuming booking.created...")
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("❌ Consumer error: %v", err)
			continue
		}

		var event BookingCreatedEvent
		if err := json.Unmarshal(m.Value, &event); err != nil {
			log.Printf("❌ Failed to decode booking.created: %v", err)
			continue
		}

		job := models.Job{
			BookingID:  event.BookingID,
			PickupLoc:  event.PickupLoc,
			Dropoff:    event.Dropoff,
			Price:      event.Price,
			RideStatus: event.RideStatus,
		}
		if err := c.repo.InsertJob(ctx, job); err != nil {
			log.Printf("❌ Failed to insert job: %v", err)
			continue
		}
		log.Printf("✅ New job inserted: booking_id=%s", event.BookingID)
	}
}
