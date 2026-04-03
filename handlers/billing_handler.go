package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"subscription-management_backend/models"
)

type CreateBillingRequest struct {
	UserID      uint     `json:"user_id" binding:"required"`
	PlanID      uint     `json:"plan_id" binding:"required"`
	Amount      *float64 `json:"amount" binding:"omitempty,gte=0"`
	DueDate     string   `json:"due_date" binding:"required"`
	Description string   `json:"description"`
}

func (h *Handler) createBilling(c *gin.Context) {
	var req CreateBillingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ensureUserExists(h.db, req.UserID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var plan models.Plan
	if err := h.db.First(&plan, req.PlanID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "plan not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dueDate, err := time.Parse(time.RFC3339, req.DueDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "due_date must be RFC3339 format"})
		return
	}

	amount := plan.Price
	if req.Amount != nil {
		amount = *req.Amount
	}

	billing := models.Billing{
		UserID:      req.UserID,
		PlanID:      req.PlanID,
		Amount:      amount,
		Status:      models.BillingStatusPending,
		DueDate:     dueDate,
		Description: req.Description,
	}

	if err := h.db.Create(&billing).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Preload("User").Preload("Plan").First(&billing, billing.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, billing)
}

func (h *Handler) listBillings(c *gin.Context) {
	var billings []models.Billing
	if err := h.db.Preload("User").Preload("Plan").Order("id DESC").Find(&billings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, billings)
}

func (h *Handler) getBilling(c *gin.Context) {
	billingID, err := parseID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid billing id"})
		return
	}

	var billing models.Billing
	if err := h.db.Preload("User").Preload("Plan").First(&billing, billingID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "billing not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, billing)
}

func (h *Handler) payBilling(c *gin.Context) {
	billingID, err := parseID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid billing id"})
		return
	}

	var billing models.Billing
	if err := h.db.First(&billing, billingID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "billing not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	now := time.Now().UTC()
	billing.Status = models.BillingStatusPaid
	billing.PaidAt = &now

	if err := h.db.Save(&billing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Preload("User").Preload("Plan").First(&billing, billing.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, billing)
}
