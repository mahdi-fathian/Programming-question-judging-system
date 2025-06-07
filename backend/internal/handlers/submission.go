package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/onlinejudge/backend/internal/models"
	"github.com/onlinejudge/backend/pkg/broker"
	"gorm.io/gorm"
)

type SubmissionHandler struct {
	db     *gorm.DB
	broker *broker.NATSClient
}

func NewSubmissionHandler(db *gorm.DB, broker *broker.NATSClient) *SubmissionHandler {
	return &SubmissionHandler{
		db:     db,
		broker: broker,
	}
}

func (h *SubmissionHandler) Submit(c *gin.Context) {
	var submission models.Submission
	if err := c.ShouldBindJSON(&submission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	submission.UserID = user.(*models.User).ID

	// Verify problem exists
	var problem models.Problem
	if err := h.db.First(&problem, submission.ProblemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "problem not found"})
		return
	}

	// Set initial status
	submission.Status = "pending"

	// Save submission
	if err := h.db.Create(&submission).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Publish to NATS for evaluation
	if err := h.broker.PublishSubmission(&submission); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to queue submission"})
		return
	}

	c.JSON(http.StatusCreated, submission)
}

func (h *SubmissionHandler) GetSubmission(c *gin.Context) {
	id := c.Param("id")
	var submission models.Submission

	if err := h.db.Preload("User").Preload("Problem").First(&submission, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "submission not found"})
		return
	}

	c.JSON(http.StatusOK, submission)
}

func (h *SubmissionHandler) ListSubmissions(c *gin.Context) {
	var submissions []models.Submission
	query := h.db.Model(&models.Submission{})

	// Apply filters
	if problemID := c.Query("problem_id"); problemID != "" {
		query = query.Where("problem_id = ?", problemID)
	}

	if userID := c.Query("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	var total int64
	query.Count(&total)

	if err := query.Preload("User").Preload("Problem").
		Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&submissions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": submissions,
		"meta": gin.H{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func (h *SubmissionHandler) GetSubmissionResults(c *gin.Context) {
	id := c.Param("id")
	var submission models.Submission

	if err := h.db.Preload("Results").First(&submission, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "submission not found"})
		return
	}

	c.JSON(http.StatusOK, submission.Results)
} 