package repository

import (
	"azh/internal/model"
	"gorm.io/gorm"
	"time"
)

// CourseRepository handles database operations for courses
type CourseRepository struct {
	db *gorm.DB
}

// NewCourseRepository creates a new CourseRepository
func NewCourseRepository(db *gorm.DB) *CourseRepository {
	return &CourseRepository{db: db}
}

// GetAll retrieves all courses, filtered by first_schedule and last_schedule within 8 days of today
func (r *CourseRepository) GetAll() ([]model.Course, error) {
	var courses []model.Course
	today := time.Now()
	eightDaysAgo := today.AddDate(0, 0, -8)
	eightDaysAhead := today.AddDate(0, 0, 8)

	err := r.db.Where("((first_schedule IS NULL OR first_schedule <= ?) AND (last_schedule IS NULL OR last_schedule >= ?))", eightDaysAhead, eightDaysAgo).
		Order("id ASC, start_time ASC").
		Find(&courses).Error
	return courses, err
}

// GetByID retrieves a course by ID
func (r *CourseRepository) GetByID(id string) (model.Course, error) {
	var course model.Course
	err := r.db.Where("id = ?", id).First(&course).Error
	return course, err
}
