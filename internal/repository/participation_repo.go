package repository

import (
	"azh/internal/model"
	"gorm.io/gorm"
	"time"
)

// ParticipationRepository handles database operations for participations
type ParticipationRepository struct {
	db *gorm.DB
}

// NewParticipationRepository creates a new ParticipationRepository
func NewParticipationRepository(db *gorm.DB) *ParticipationRepository {
	return &ParticipationRepository{db: db}
}

// GetByCourseAndDate retrieves participations for a specific course and date
func (r *ParticipationRepository) GetByCourseAndDate(courseID, date string) ([]model.Participation, error) {
	var participations []model.Participation
	err := r.db.Where("course_id = ? AND date = ?", courseID, date).Find(&participations).Error
	return participations, err
}

// Upsert updates or inserts a participation record
func (r *ParticipationRepository) Upsert(participation *model.Participation) error {
	return r.db.Save(participation).Error
}

// Delete removes a participation record
func (r *ParticipationRepository) Delete(memberID, courseID uint, date time.Time) error {
	return r.db.Unscoped().Where("member_id = ? AND course_id = ? AND date = ?", memberID, courseID, date).Delete(&model.Participation{}).Error
}

// GetExportData retrieves participation data within a date range for export
func (r *ParticipationRepository) GetExportData(minDate, maxDate string) ([]model.Participation, error) {
	var participations []model.Participation
	err := r.db.Where("date >= ? AND date <= ?", minDate, maxDate).
		Order("date ASC, course_id ASC, member_id ASC").
		Find(&participations).Error
	return participations, err
}
