package handlers

import (
	"context"
	"driver_svc/internal/mq"
	"driver_svc/internal/repo"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type DriverHandler struct {
	Repo     *repo.DriverRepo
	Producer *mq.Producer
}

func NewDriverHandler(r *repo.DriverRepo, p *mq.Producer) *DriverHandler {
	return &DriverHandler{Repo: r, Producer: p}
}

// ✅ GET /drivers
func (h *DriverHandler) GetDrivers(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	drivers, err := h.Repo.GetDrivers(ctx)
	if err != nil {
		http.Error(w, "failed to fetch drivers", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(drivers)
}

// ✅ GET /jobs
func (h *DriverHandler) GetJobs(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	jobs, err := h.Repo.GetJobs(ctx)
	if err != nil {
		http.Error(w, "failed to fetch jobs", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(jobs)
}

// ✅ POST /jobs/{booking_id}/accept
func (h *DriverHandler) AcceptJob(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	bookingID := chi.URLParam(r, "booking_id")

	var req struct {
		DriverID string `json:"driver_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	job, err := h.Repo.AcceptJob(ctx, bookingID, req.DriverID)
	if err != nil {
		http.Error(w, "job already accepted", http.StatusConflict)
		return
	}

	// Produce booking.accepted event
	_ = h.Producer.PublishBookingAccepted(ctx, bookingID, req.DriverID)

	json.NewEncoder(w).Encode(job)
}
