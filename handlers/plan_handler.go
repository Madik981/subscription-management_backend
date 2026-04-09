package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"subscription-management_backend/models"
)

type Handler struct {
	db        *gorm.DB
	jwtSecret []byte
	tokenTTL  time.Duration
}

type CreatePlanRequest struct {
	Name         string  `json:"name" binding:"required"`
	Description  string  `json:"description"`
	Price        float64 `json:"price" binding:"required,gte=0"`
	Currency     string  `json:"currency"`
	BillingCycle string  `json:"billing_cycle"`
}

type UpdatePlanRequest struct {
	Name         *string  `json:"name"`
	Description  *string  `json:"description"`
	Price        *float64 `json:"price" binding:"omitempty,gte=0"`
	Currency     *string  `json:"currency"`
	BillingCycle *string  `json:"billing_cycle"`
}

func NewHandler(db *gorm.DB, jwtSecret string) *Handler {
	return &Handler{
		db:        db,
		jwtSecret: []byte(jwtSecret),
		tokenTTL:  24 * time.Hour,
	}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	router.GET("/health", h.health)
	auth := router.Group("/auth")
	{
		auth.POST("/register", h.register)
		auth.POST("/login", h.login)
		auth.GET("/me", h.authMiddleware(), h.me)
	}

	protected := router.Group("/")
	protected.Use(h.authMiddleware())

	plans := protected.Group("/plans")
	{
		plans.POST("", h.createPlan)
		plans.GET("", h.listPlans)
		plans.GET("/:id", h.getPlan)
		plans.PATCH("/:id", h.updatePlan)
	}

	users := protected.Group("/users")
	{
		users.POST("", h.createUser)
		users.GET("", h.listUsers)
		users.GET("/:id", h.getUser)
		users.PATCH("/:id", h.updateUser)
	}

	billings := protected.Group("/billings")
	{
		billings.POST("", h.createBilling)
		billings.GET("", h.listBillings)
		billings.GET("/:id", h.getBilling)
		billings.PATCH("/:id/pay", h.payBilling)
	}
}

func (h *Handler) health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handler) createPlan(c *gin.Context) {
	var req CreatePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plan := models.Plan{
		Name:         req.Name,
		Description:  req.Description,
		Price:        req.Price,
		Currency:     coalesceString(req.Currency, "USD"),
		BillingCycle: coalesceString(req.BillingCycle, "monthly"),
	}

	if err := h.db.Create(&plan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, plan)
}

func (h *Handler) listPlans(c *gin.Context) {
	var plans []models.Plan
	if err := h.db.Order("id DESC").Find(&plans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, plans)
}

func (h *Handler) getPlan(c *gin.Context) {
	planID, err := parseID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan id"})
		return
	}

	var plan models.Plan
	if err := h.db.First(&plan, planID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "plan not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, plan)
}

func (h *Handler) updatePlan(c *gin.Context) {
	planID, err := parseID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan id"})
		return
	}

	var req UpdatePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var plan models.Plan
	if err := h.db.First(&plan, planID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "plan not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if req.Name != nil {
		plan.Name = *req.Name
	}
	if req.Description != nil {
		plan.Description = *req.Description
	}
	if req.Price != nil {
		plan.Price = *req.Price
	}
	if req.Currency != nil {
		plan.Currency = *req.Currency
	}
	if req.BillingCycle != nil {
		plan.BillingCycle = *req.BillingCycle
	}

	if err := h.db.Save(&plan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, plan)
}

func parseID(rawID string) (uint, error) {
	id, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

func coalesceString(value string, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
