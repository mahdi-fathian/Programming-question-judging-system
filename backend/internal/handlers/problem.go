package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/onlinejudge/backend/internal/models"
	"github.com/onlinejudge/backend/pkg/database"
)

type CreateProblemRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Difficulty  string `json:"difficulty" binding:"required"`
	TimeLimit   int    `json:"time_limit" binding:"required"`
	MemoryLimit int    `json:"memory_limit" binding:"required"`
	TestCases   []struct {
		Input    string `json:"input" binding:"required"`
		Output   string `json:"output" binding:"required"`
		IsSample bool   `json:"is_sample"`
	} `json:"test_cases" binding:"required,min=1"`
	Tags []string `json:"tags"`
}

func CreateProblem(c *gin.Context) {
	var req CreateProblemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, _ := c.Get("user")
	u := user.(models.User)

	// Create problem
	problem := models.Problem{
		Title:       req.Title,
		Description: req.Description,
		Difficulty:  req.Difficulty,
		TimeLimit:   req.TimeLimit,
		MemoryLimit: req.MemoryLimit,
		CreatedBy:   u.ID,
	}

	// Create test cases
	for _, tc := range req.TestCases {
		problem.TestCases = append(problem.TestCases, models.TestCase{
			Input:    tc.Input,
			Output:   tc.Output,
			IsSample: tc.IsSample,
		})
	}

	// Create or get tags
	for _, tagName := range req.Tags {
		var tag models.Tag
		database.DB.FirstOrCreate(&tag, models.Tag{Name: tagName})
		problem.Tags = append(problem.Tags, tag)
	}

	if err := database.DB.Create(&problem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create problem"})
		return
	}

	c.JSON(http.StatusCreated, problem)
}

func GetProblem(c *gin.Context) {
	id := c.Param("id")
	var problem models.Problem

	if err := database.DB.Preload("TestCases").Preload("Tags").First(&problem, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Problem not found"})
		return
	}

	c.JSON(http.StatusOK, problem)
}

func ListProblems(c *gin.Context) {
	var problems []models.Problem
	query := database.DB.Model(&models.Problem{})

	// Apply filters
	if difficulty := c.Query("difficulty"); difficulty != "" {
		query = query.Where("difficulty = ?", difficulty)
	}

	if tag := c.Query("tag"); tag != "" {
		query = query.Joins("JOIN problem_tags ON problem_tags.problem_id = problems.id").
			Joins("JOIN tags ON tags.id = problem_tags.tag_id").
			Where("tags.name = ?", tag)
	}

	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var total int64
	query.Count(&total)

	if err := query.Offset(offset).Limit(limit).Find(&problems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch problems"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"problems": problems,
		"total":    total,
		"page":     page,
		"limit":    limit,
	})
}

func UpdateProblem(c *gin.Context) {
	id := c.Param("id")
	var problem models.Problem

	if err := database.DB.First(&problem, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Problem not found"})
		return
	}

	// Check if user is the creator or an admin
	user, _ := c.Get("user")
	u := user.(models.User)
	if problem.CreatedBy != u.ID && u.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to update this problem"})
		return
	}

	var updateData CreateProblemRequest
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update problem fields
	problem.Title = updateData.Title
	problem.Description = updateData.Description
	problem.Difficulty = updateData.Difficulty
	problem.TimeLimit = updateData.TimeLimit
	problem.MemoryLimit = updateData.MemoryLimit

	// Update test cases
	database.DB.Where("problem_id = ?", problem.ID).Delete(&models.TestCase{})
	for _, tc := range updateData.TestCases {
		problem.TestCases = append(problem.TestCases, models.TestCase{
			Input:    tc.Input,
			Output:   tc.Output,
			IsSample: tc.IsSample,
		})
	}

	// Update tags
	database.DB.Model(&problem).Association("Tags").Clear()
	for _, tagName := range updateData.Tags {
		var tag models.Tag
		database.DB.FirstOrCreate(&tag, models.Tag{Name: tagName})
		problem.Tags = append(problem.Tags, tag)
	}

	if err := database.DB.Save(&problem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update problem"})
		return
	}

	c.JSON(http.StatusOK, problem)
}

func DeleteProblem(c *gin.Context) {
	id := c.Param("id")
	var problem models.Problem

	if err := database.DB.First(&problem, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Problem not found"})
		return
	}

	// Check if user is the creator or an admin
	user, _ := c.Get("user")
	u := user.(models.User)
	if problem.CreatedBy != u.ID && u.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to delete this problem"})
		return
	}

	if err := database.DB.Delete(&problem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete problem"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Problem deleted successfully"})
} 