package handlers

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"gorm.io/gorm"
)

type Handler struct {
	db             *gorm.DB
	jwtSecret      []byte
	tokenTTL       time.Duration
	billingBaseURL string
	internalToken  string
	httpClient     *resty.Client
}

func NewHandler(db *gorm.DB, jwtSecret string, billingBaseURL string, internalToken string) *Handler {
	trimmedBillingURL := strings.TrimRight(billingBaseURL, "/")
	client := resty.New().
		SetTimeout(5*time.Second).
		SetHeader("X-Internal-Token", internalToken)

	return &Handler{
		db:             db,
		jwtSecret:      []byte(jwtSecret),
		tokenTTL:       24 * time.Hour,
		billingBaseURL: trimmedBillingURL,
		internalToken:  internalToken,
		httpClient:     client,
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

	users := protected.Group("/users")
	{
		users.POST("", h.createUser)
		users.GET("", h.listUsers)
		users.GET("/:id", h.getUser)
		users.PATCH("/:id", h.updateUser)
	}

	internal := router.Group("/internal")
	internal.Use(h.internalAuthMiddleware())
	{
		internal.GET("/users/:id", h.internalGetUser)
	}
}

func (h *Handler) health(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}
