package model

import (
	"gorm.io/gorm"
	"time"
)

// Participation represents a member's attendance at a specific course on a specific date
type Participation struct {
	gorm.Model
	MemberID uint      `gorm:"index" json:"member_id"`
	CourseID uint      `gorm:"index" json:"course_id"`
	Date     time.Time `gorm:"index;type:date" json:"date"` // Format: YYYY-MM-DD
}
