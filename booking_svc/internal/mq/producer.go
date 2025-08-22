package mq

import (
	"booking_svc/internal/models"
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
	topic  string
}

// ✅ Initialize Kafka producer
func NewProducer(broker string, topic string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(broker),
			Balancer: &kafka.LeastBytes{},
		},
		topic: topic,
	}
}

// ✅ Publish booking.created event
func (p *Producer) PublishBookingCreated(ctx context.Context, booking models.Booking) error {
	payload, err := json.Marshal(map[string]interface{}{
		"booking_id":  booking.ID.Hex(),
		"pickuploc":   booking.PickupLoc,
		"dropoff":     booking.Dropoff,
		"price":       booking.Price,
		"ride_status": booking.RideStatus,
	})
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Key:   []byte(booking.ID.Hex()),
		Value: payload,
		Topic: p.topic,
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		log.Printf("❌ Failed to publish booking.created: %v", err)
		return err
	}

	log.Printf("✅ Published booking.created for booking_id=%s", booking.ID.Hex())
	return nil
}

// ✅ Close producer (for graceful shutdown)
func (p *Producer) Close() error {
	return p.writer.Close()
}
