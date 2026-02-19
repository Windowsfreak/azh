package model

import "gorm.io/gorm"

// MemberCourse represents the n:m relationship between members and courses
type MemberCourse struct {
	gorm.Model
	MemberID uint `gorm:"index" json:"member_id"`
	CourseID uint `gorm:"index" json:"course_id"`
}
