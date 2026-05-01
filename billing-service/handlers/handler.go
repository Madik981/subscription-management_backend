package handlers

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"gorm.io/gorm"
)

type Handler struct {
	db              *gorm.DB
	jwtSecret       []byte
	accountsBaseURL string
	internalToken   string
	httpClient      *resty.Client
}

func NewHandler(db *gorm.DB, jwtSecret string, accountsBaseURL string, internalToken string) *Handler {
	trimmedAccountsURL := strings.TrimRight(accountsBaseURL, "/")
	client := resty.New().
		SetTimeout(5*time.Second).
		SetHeader("X-Internal-Token", internalToken)

	return &Handler{
		db:              db,
		jwtSecret:       []byte(jwtSecret),
		accountsBaseURL: trimmedAccountsURL,
		internalToken:   internalToken,
		httpClient:      client,
	}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	router.GET("/health", h.health)

	internal := router.Group("/internal")
	internal.Use(h.internalAuthMiddleware())
	{
		internal.GET("/plans/:id", h.internalGetPlan)
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

	billings := protected.Group("/billings")
	{
		billings.POST("", h.createBilling)
		billings.GET("", h.listBillings)
		billings.GET("/:id", h.getBilling)
		billings.PATCH("/:id/pay", h.payBilling)
	}
}

func (h *Handler) health(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}
