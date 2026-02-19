package repository

import (
	"azh/internal/model"
	"gorm.io/gorm"
	"time"
)

// MemberRepository handles database operations for members
type MemberRepository struct {
	db *gorm.DB
}

// NewMemberRepository creates a new MemberRepository
func NewMemberRepository(db *gorm.DB) *MemberRepository {
	return &MemberRepository{db: db}
}

// GetByIDs retrieves members by their IDs
func (r *MemberRepository) GetByIDs(ids []uint) ([]model.Member, error) {
	var members []model.Member
	err := r.db.Where("id IN ?", ids).
		Order("first_name ASC, last_name ASC, id ASC").
		Find(&members).Error
	return members, err
}

// GetByIDsAndDate retrieves members by their IDs, filtered by sign_up_date and cancellation_date for a specific date
func (r *MemberRepository) GetByIDsAndDate(ids []uint, selectedDate time.Time) ([]model.Member, error) {
	var members []model.Member
	err := r.db.Where("id IN ?", ids).
		Where("(sign_up_date IS NULL OR sign_up_date <= ?)", selectedDate).
		Where("(cancellation_date IS NULL OR cancellation_date >= ?)", selectedDate).
		Order("first_name ASC, last_name ASC, id ASC").
		Find(&members).Error
	return members, err
}
