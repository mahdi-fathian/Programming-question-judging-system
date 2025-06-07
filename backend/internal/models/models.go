package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;not null"`
	Email    string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
	Role     string `gorm:"default:user"`
}

type Problem struct {
	gorm.Model
	Title       string `gorm:"not null"`
	Description string `gorm:"type:text;not null"`
	Difficulty  string `gorm:"not null"`
	Points      int    `gorm:"not null"`
	TimeLimit   int    `gorm:"not null"` // in milliseconds
	MemoryLimit int    `gorm:"not null"` // in MB
	TestCases   []TestCase
	Submissions []Submission
	Contests    []Contest `gorm:"many2many:contest_problems;"`
}

type TestCase struct {
	gorm.Model
	ProblemID  uint   `gorm:"not null"`
	Input      string `gorm:"type:text;not null"`
	Output     string `gorm:"type:text;not null"`
	IsExample  bool   `gorm:"default:false"`
	Problem    Problem
}

type Contest struct {
	gorm.Model
	Title       string    `gorm:"not null"`
	Description string    `gorm:"type:text;not null"`
	StartTime   time.Time `gorm:"not null"`
	EndTime     time.Time `gorm:"not null"`
	IsPublic    bool      `gorm:"default:true"`
	CreatorID   uint      `gorm:"not null"`
	Creator     User      `gorm:"foreignKey:CreatorID"`
	Problems    []Problem `gorm:"many2many:contest_problems;"`
	Users       []User    `gorm:"many2many:contest_users;"`
	Submissions []Submission
}

type Submission struct {
	gorm.Model
	UserID    uint   `gorm:"not null"`
	User      User   `gorm:"foreignKey:UserID"`
	ProblemID uint   `gorm:"not null"`
	Problem   Problem `gorm:"foreignKey:ProblemID"`
	ContestID *uint
	Contest   *Contest `gorm:"foreignKey:ContestID"`
	Language  string   `gorm:"not null"`
	Code      string   `gorm:"type:text;not null"`
	Status    string   `gorm:"not null"`
	Score     int      `gorm:"default:0"`
	Results   []SubmissionResult
}

type SubmissionResult struct {
	gorm.Model
	SubmissionID uint      `gorm:"not null"`
	Submission   Submission `gorm:"foreignKey:SubmissionID"`
	TestCaseID   uint      `gorm:"not null"`
	TestCase     TestCase  `gorm:"foreignKey:TestCaseID"`
	Status       string    `gorm:"not null"`
	TimeUsed     int       `gorm:"not null"` // in milliseconds
	MemoryUsed   int       `gorm:"not null"` // in KB
	Error        string    `gorm:"type:text"`
} 