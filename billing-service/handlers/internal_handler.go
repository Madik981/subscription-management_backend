package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"subscription-management/billing-service/models"
)

func (h *Handler) internalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if h.internalToken == "" {
			c.Next()
			return
		}

		if c.GetHeader("X-Internal-Token") != h.internalToken {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid internal token"})
			return
		}

		c.Next()
	}
}

func (h *Handler) internalGetPlan(c *gin.Context) {
	planID, err := parseID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan id"})
		return
	}

	var plan models.Plan
	if err := h.db.Select("id").First(&plan, planID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "plan not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": plan.ID})
}
