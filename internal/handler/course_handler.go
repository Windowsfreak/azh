package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"azh/internal/service"
	"github.com/julienschmidt/httprouter"
)

// CourseHandler handles HTTP requests for courses
type CourseHandler struct {
	courseService *service.CourseService
}

// NewCourseHandler creates a new CourseHandler
func NewCourseHandler(courseService *service.CourseService) *CourseHandler {
	return &CourseHandler{courseService: courseService}
}

// GetCourses handles GET /api/courses
func (h *CourseHandler) GetCourses(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	courses, err := h.courseService.GetCourses()
	if err != nil {
		http.Error(w, "Failed to retrieve courses", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

// GetOccurrences handles GET /api/courses/:id/occurrences
func (h *CourseHandler) GetOccurrences(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	courseID := ps.ByName("id")
	referenceDate := time.Now() // Use current date as reference; can be parameterized if needed
	occurrences, err := h.courseService.GetOccurrences(courseID, referenceDate)
	if err != nil {
		http.Error(w, "Failed to calculate occurrences", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(occurrences)
}
