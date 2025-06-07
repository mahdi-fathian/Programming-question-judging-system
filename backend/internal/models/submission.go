package models

import (
	"time"
)

type Submission struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	ProblemID uint      `json:"problem_id" gorm:"not null"`
	ContestID *uint     `json:"contest_id"` // Optional, nil if not part of a contest
	Language  string    `json:"language" gorm:"not null"` // e.g., "cpp", "python", "java"
	Code      string    `json:"code" gorm:"type:text;not null"`
	Status    string    `json:"status" gorm:"not null"` // pending, accepted, wrong_answer, time_limit, memory_limit, runtime_error, compilation_error
	TimeUsed  int       `json:"time_used"` // in milliseconds
	MemoryUsed int      `json:"memory_used"` // in KB
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	User    User    `json:"user" gorm:"foreignKey:UserID"`
	Problem Problem `json:"problem" gorm:"foreignKey:ProblemID"`
	Contest *Contest `json:"contest,omitempty" gorm:"foreignKey:ContestID"`
}

type SubmissionResult struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	SubmissionID uint      `json:"submission_id" gorm:"not null"`
	TestCaseID   uint      `json:"test_case_id" gorm:"not null"`
	Status       string    `json:"status" gorm:"not null"` // accepted, wrong_answer, time_limit, memory_limit, runtime_error
	TimeUsed     int       `json:"time_used"` // in milliseconds
	MemoryUsed   int       `json:"memory_used"` // in KB
	Error        string    `json:"error" gorm:"type:text"` // error message if any
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relationships
	Submission Submission `json:"submission" gorm:"foreignKey:SubmissionID"`
	TestCase   TestCase   `json:"test_case" gorm:"foreignKey:TestCaseID"`
} 