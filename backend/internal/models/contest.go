package models

import (
	"time"
)

type Contest struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description" gorm:"type:text"`
	StartTime   time.Time `json:"start_time" gorm:"not null"`
	EndTime     time.Time `json:"end_time" gorm:"not null"`
	CreatedBy   uint      `json:"created_by" gorm:"not null"`
	IsPublic    bool      `json:"is_public" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	Problems []Problem `json:"problems" gorm:"many2many:contest_problems;"`
	Users    []User    `json:"users" gorm:"many2many:contest_users;"`
}

type ContestProblem struct {
	ContestID  uint `json:"contest_id" gorm:"primaryKey"`
	ProblemID  uint `json:"problem_id" gorm:"primaryKey"`
	Order      int  `json:"order" gorm:"not null"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ContestUser struct {
	ContestID  uint      `json:"contest_id" gorm:"primaryKey"`
	UserID     uint      `json:"user_id" gorm:"primaryKey"`
	JoinedAt   time.Time `json:"joined_at"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
} 