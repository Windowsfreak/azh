package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"azh/internal/service"
	"github.com/julienschmidt/httprouter"
)

// ParticipationHandler handles HTTP requests for participations
type ParticipationHandler struct {
	participationService *service.ParticipationService
}

// NewParticipationHandler creates a new ParticipationHandler
func NewParticipationHandler(participationService *service.ParticipationService) *ParticipationHandler {
	return &ParticipationHandler{participationService: participationService}
}

// GetParticipants handles GET /api/courses/:id/dates/:date/participants
func (h *ParticipationHandler) GetParticipants(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	courseID := ps.ByName("id")
	date := ps.ByName("date")
	participants, err := h.participationService.GetParticipants(courseID, date)
	if err != nil {
		http.Error(w, "Failed to retrieve participants", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(participants)
}

// SetAttendance handles POST /api/courses/:id/dates/:date/participants/:participantId/attendance
func (h *ParticipationHandler) SetAttendance(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	courseIDStr := ps.ByName("id")
	courseID, err := strconv.ParseUint(courseIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid course ID", http.StatusBadRequest)
		return
	}
	dateStr := ps.ByName("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}
	participantIDStr := ps.ByName("participantId")
	participantID, err := strconv.ParseUint(participantIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid participant ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Present bool `json:"present"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.participationService.SetAttendance(uint(courseID), date, uint(participantID), req.Present)
	if err != nil {
		http.Error(w, "Failed to update attendance", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"present": req.Present})
}

// ExportData handles GET /api/export?minDate=YYYY-MM-DD&maxDate=YYYY-MM-DD
func (h *ParticipationHandler) ExportData(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	query := r.URL.Query()
	minDate := query.Get("minDate")
	maxDate := query.Get("maxDate")
	if minDate == "" || maxDate == "" {
		http.Error(w, "minDate and maxDate are required", http.StatusBadRequest)
		return
	}

	csvData, err := h.participationService.ExportData(minDate, maxDate)
	if err != nil {
		http.Error(w, "Failed to export data", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=participations.csv")
	w.Write([]byte(csvData))
}
