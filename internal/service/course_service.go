package service

import (
	"azh/internal/model"
	"azh/internal/repository"
	"time"
)

// CourseService handles business logic for courses
type CourseService struct {
	courseRepo *repository.CourseRepository
}

// NewCourseService creates a new CourseService
func NewCourseService(courseRepo *repository.CourseRepository) *CourseService {
	return &CourseService{courseRepo: courseRepo}
}

// GetCourses retrieves all courses, filtered by first_schedule and last_schedule within 8 days of today
func (s *CourseService) GetCourses() ([]model.Course, error) {
	return s.courseRepo.GetAll()
}

// GetOccurrences calculates the previous, current, and next occurrence of a course
func (s *CourseService) GetOccurrences(courseID string, referenceDate time.Time) ([]string, error) {
	course, err := s.courseRepo.GetByID(courseID)
	if err != nil {
		return nil, err
	}
	return calculateOccurrences(course.Weekday, referenceDate), nil
}

// calculateOccurrences computes the three occurrences based on weekday
func calculateOccurrences(courseWeekday string, referenceDate time.Time) []string {
	weekdayMap := map[string]time.Weekday{
		"Montag":     time.Monday,
		"Dienstag":   time.Tuesday,
		"Mittwoch":   time.Wednesday,
		"Donnerstag": time.Thursday,
		"Freitag":    time.Friday,
		"Samstag":    time.Saturday,
		"Sonntag":    time.Sunday,
	}
	courseWeekdayInt := weekdayMap[courseWeekday]
	currentWeekdayInt := referenceDate.Weekday()
	daysToAdd := int(courseWeekdayInt) - int(currentWeekdayInt)
	if daysToAdd < 0 {
		daysToAdd += 7
	}
	current := referenceDate.AddDate(0, 0, daysToAdd)
	previous := current.AddDate(0, 0, -7)
	next := current.AddDate(0, 0, 7)
	return []string{
		previous.Format("2006-01-02"),
		current.Format("2006-01-02"),
		next.Format("2006-01-02"),
	}
}
