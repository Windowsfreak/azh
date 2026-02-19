package repository

import (
	"azh/internal/model"
	"gorm.io/gorm"
)

// MemberCourseRepository handles database operations for member-course relationships
type MemberCourseRepository struct {
	db *gorm.DB
}

// NewMemberCourseRepository creates a new MemberCourseRepository
func NewMemberCourseRepository(db *gorm.DB) *MemberCourseRepository {
	return &MemberCourseRepository{db: db}
}

// GetMembersByCourseID retrieves member IDs associated with a course
func (r *MemberCourseRepository) GetMembersByCourseID(courseID string) ([]uint, error) {
	var memberCourses []model.MemberCourse
	err := r.db.Where("course_id = ?", courseID).Find(&memberCourses).Error
	if err != nil {
		return nil, err
	}
	memberIDs := make([]uint, len(memberCourses))
	for i, mc := range memberCourses {
		memberIDs[i] = mc.MemberID
	}
	return memberIDs, nil
}
