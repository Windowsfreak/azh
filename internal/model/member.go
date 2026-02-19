package model

import (
	"gorm.io/gorm"
	"time"
)

// Member represents a club member
type Member struct {
	gorm.Model
	ID               uint      `gorm:"primaryKey" json:"id"`
	FirstName        string    `gorm:"type:varchar(100)" json:"first_name"`
	LastName         string    `gorm:"type:varchar(100)" json:"last_name"`
	Email            string    `gorm:"type:varchar(255)" json:"email"`
	Phone            string    `gorm:"type:varchar(50)" json:"phone"`
	SignUpDate       time.Time `gorm:"type:date" json:"sign_up_date"`
	CancellationDate time.Time `gorm:"type:date" json:"cancellation_date"`
	Age              int       `json:"age"`
	Notes            string    `gorm:"type:text" json:"notes"`
}
