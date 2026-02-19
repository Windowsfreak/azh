package model

import (
	"gorm.io/gorm"
	"time"
)

// Course represents a training course
type Course struct {
	gorm.Model
	ID            uint      `gorm:"primaryKey" json:"id"`
	Name          string    `gorm:"type:varchar(255)" json:"name"`
	Location      string    `gorm:"type:varchar(100)" json:"location"`
	TrainingType  string    `gorm:"type:varchar(100)" json:"training_type"`
	Weekday       string    `gorm:"type:varchar(20)" json:"weekday"`
	StartTime     string    `gorm:"type:varchar(10)" json:"start_time"`
	EndTime       string    `gorm:"type:varchar(10)" json:"end_time"`
	FirstSchedule time.Time `gorm:"type:date" json:"first_schedule"`
	LastSchedule  time.Time `gorm:"type:date" json:"last_schedule"`
	TrainerNames  string    `gorm:"type:text" json:"trainer_names"`
}
