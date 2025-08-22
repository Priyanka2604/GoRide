package mq

import (
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

type Producer struct {
	writer *kafka.Writer
	topic  string
}

func NewProducer(broker, topic string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(broker),
			Balancer: &kafka.LeastBytes{},
		},
		topic: topic,
	}
}

func (p *Producer) PublishBookingAccepted(ctx context.Context, bookingID, driverID string) error {
	event := AcceptedEvent{
		BookingID:  bookingID,
		DriverID:   driverID,
		RideStatus: "Accepted",
	}

	payload, _ := json.Marshal(event)

	err := p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(bookingID),
		Value: payload,
		Topic: p.topic,
	})
	if err != nil {
		log.Printf("❌ Failed to publish booking.accepted: %v", err)
		return err
	}
	log.Printf("✅ Published booking.accepted for booking_id=%s driver=%s", bookingID, driverID)
	return nil
}
