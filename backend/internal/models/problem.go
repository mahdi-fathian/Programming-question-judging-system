package models

import (
	"time"
)

type Problem struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description" gorm:"type:text;not null"`
	Difficulty  string    `json:"difficulty" gorm:"not null"` // easy, medium, hard
	TimeLimit   int       `json:"time_limit" gorm:"not null"` // in milliseconds
	MemoryLimit int       `json:"memory_limit" gorm:"not null"` // in MB
	CreatedBy   uint      `json:"created_by" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	// Relationships
	TestCases []TestCase `json:"test_cases" gorm:"foreignKey:ProblemID"`
	Tags      []Tag      `json:"tags" gorm:"many2many:problem_tags;"`
}

type TestCase struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ProblemID uint      `json:"problem_id" gorm:"not null"`
	Input     string    `json:"input" gorm:"type:text;not null"`
	Output    string    `json:"output" gorm:"type:text;not null"`
	IsSample  bool      `json:"is_sample" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Tag struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"unique;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
} 