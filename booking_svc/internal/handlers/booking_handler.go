package handlers

import (
	"booking_svc/internal/models"
	"booking_svc/internal/mq"
	"booking_svc/internal/repo"
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type BookingHandler struct {
	Repo     *repo.BookingRepo
	Producer *mq.Producer
}

func NewBookingHandler(r *repo.BookingRepo, p *mq.Producer) *BookingHandler {
	return &BookingHandler{Repo: r, Producer: p}
}

// ✅ POST /bookings
func (h *BookingHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var req struct {
		PickupLoc models.Location `json:"pickuploc"`
		Dropoff   models.Location `json:"dropoff"`
		Price     int             `json:"price"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	booking := models.Booking{
		PickupLoc:  req.PickupLoc,
		Dropoff:    req.Dropoff,
		Price:      req.Price,
		RideStatus: "Requested",
		DriverID:   nil,
		CreatedAt:  models.NowUTC(),
	}

	created, err := h.Repo.CreateBooking(ctx, booking)
	if err != nil {
		http.Error(w, "failed to create booking", http.StatusInternalServerError)
		return
	}

	// Publish booking.created event
	_ = h.Producer.PublishBookingCreated(ctx, *created)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

// ✅ GET /bookings
func (h *BookingHandler) GetAllBookings(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	bookings, err := h.Repo.GetAllBookings(ctx)
	if err != nil {
		http.Error(w, "failed to fetch bookings", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookings)
}

// ✅ GET /bookings/:id
func (h *BookingHandler) GetBookingByID(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	id := chi.URLParam(r, "id")

	booking, err := h.Repo.GetBookingByID(ctx, id)
	if err != nil {
		http.Error(w, "failed to fetch booking", http.StatusInternalServerError)
		return
	}
	if booking == nil {
		http.Error(w, "booking not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(booking)
}
