package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"backend/internal/models"
	"backend/internal/database"
)

type CreateContestRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description" binding:"required"`
	StartTime   time.Time `json:"start_time" binding:"required"`
	EndTime     time.Time `json:"end_time" binding:"required"`
	IsPublic    bool      `json:"is_public"`
	ProblemIDs  []uint    `json:"problem_ids" binding:"required,min=1"`
}

func CreateContest(c *gin.Context) {
	var req CreateContestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, _ := c.Get("user")
	u := user.(models.User)

	contest := models.Contest{
		Title:       req.Title,
		Description: req.Description,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		IsPublic:    req.IsPublic,
		CreatorID:   u.ID,
	}

	// Associate problems
	var problems []models.Problem
	if err := database.DB.Where("id IN ?", req.ProblemIDs).Find(&problems).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid problem IDs"})
		return
	}
	contest.Problems = problems

	if err := database.DB.Create(&contest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create contest"})
		return
	}

	c.JSON(http.StatusCreated, contest)
}

func GetContest(c *gin.Context) {
	id := c.Param("id")
	var contest models.Contest
	if err := database.DB.Preload("Problems").Preload("Users").First(&contest, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contest not found"})
		return
	}
	c.JSON(http.StatusOK, contest)
}

func ListContests(c *gin.Context) {
	var contests []models.Contest
	query := database.DB.Model(&models.Contest{})

	if isPublic := c.Query("is_public"); isPublic != "" {
		if isPublic == "true" {
			query = query.Where("is_public = ?", true)
		} else if isPublic == "false" {
			query = query.Where("is_public = ?", false)
		}
	}

	if now := c.Query("now"); now == "true" {
		nowTime := time.Now()
		query = query.Where("start_time <= ? AND end_time >= ?", nowTime, nowTime)
	}

	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var total int64
	query.Count(&total)

	if err := query.Preload("Problems").Offset(offset).Limit(limit).Find(&contests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch contests"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"contests": contests,
		"total":    total,
		"page":     page,
		"limit":    limit,
	})
}

func RegisterForContest(c *gin.Context) {
	id := c.Param("id")
	var contest models.Contest
	if err := database.DB.First(&contest, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contest not found"})
		return
	}

	if !contest.IsPublic {
		c.JSON(http.StatusForbidden, gin.H{"error": "Registration is not open for this contest"})
		return
	}

	if time.Now().After(contest.EndTime) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Contest has already ended"})
		return
	}

	user, _ := c.Get("user")
	u := user.(models.User)

	if err := database.DB.Model(&contest).Association("Users").Append(&u); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register for contest"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registered successfully"})
}

func UpdateContest(c *gin.Context) {
	id := c.Param("id")
	var contest models.Contest
	if err := database.DB.Preload("Problems").First(&contest, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contest not found"})
		return
	}

	user, _ := c.Get("user")
	u := user.(models.User)
	if contest.CreatorID != u.ID && u.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to update this contest"})
		return
	}

	var req CreateContestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contest.Title = req.Title
	contest.Description = req.Description
	contest.StartTime = req.StartTime
	contest.EndTime = req.EndTime
	contest.IsPublic = req.IsPublic

	// Update problems
	var problems []models.Problem
	if err := database.DB.Where("id IN ?", req.ProblemIDs).Find(&problems).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid problem IDs"})
		return
	}
	database.DB.Model(&contest).Association("Problems").Replace(problems)

	if err := database.DB.Save(&contest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update contest"})
		return
	}

	c.JSON(http.StatusOK, contest)
}

func DeleteContest(c *gin.Context) {
	id := c.Param("id")
	var contest models.Contest
	if err := database.DB.First(&contest, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contest not found"})
		return
	}

	user, _ := c.Get("user")
	u := user.(models.User)
	if contest.CreatorID != u.ID && u.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to delete this contest"})
		return
	}

	if err := database.DB.Delete(&contest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete contest"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contest deleted successfully"})
} 