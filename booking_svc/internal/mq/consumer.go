package mq

import (
	"booking_svc/internal/repo"
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

type AcceptedEvent struct {
	BookingID  string `json:"booking_id"`
	DriverID   string `json:"driver_id"`
	RideStatus string `json:"ride_status"`
}

type Consumer struct {
	reader *kafka.Reader
	repo   *repo.BookingRepo
}

func NewConsumer(broker, topic, groupID string, r *repo.BookingRepo) *Consumer {
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
	log.Println("📥 booking_svc consumer listening on topic booking.accepted...")
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("❌ Error reading message: %v", err)
			continue
		}

		var event AcceptedEvent
		if err := json.Unmarshal(m.Value, &event); err != nil {
			log.Printf("❌ Failed to parse booking.accepted event: %v", err)
			continue
		}

		// Only handle Accepted events
		if event.RideStatus != "Accepted" {
			continue
		}

		err = c.repo.UpdateBookingAccepted(ctx, event.BookingID, event.DriverID)
		if err != nil {
			log.Printf("❌ Failed to update booking %s: %v", event.BookingID, err)
			continue
		}

		log.Printf("✅ Booking %s accepted by driver %s", event.BookingID, event.DriverID)
	}
}
