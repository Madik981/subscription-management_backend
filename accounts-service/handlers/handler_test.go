package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"subscription-management/accounts-service/models"
)

const testJWTSecret = "test-secret"

func newTestRouter(t *testing.T) (*gin.Engine, *gorm.DB) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	dbName := strings.NewReplacer("/", "_", " ", "_").Replace(t.Name())
	db, err := gorm.Open(sqlite.Open("file:"+dbName+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("unwrap test db: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)
	if err := db.AutoMigrate(&models.User{}); err != nil {
		t.Fatalf("migrate test db: %v", err)
	}

	router := gin.New()
	NewHandler(db, testJWTSecret, "", "internal-token").RegisterRoutes(router)
	return router, db
}

func performJSON(router http.Handler, method, path string, body any, token string) *httptest.ResponseRecorder {
	payload, _ := json.Marshal(body)
	req := httptest.NewRequest(method, path, bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func registerAndToken(t *testing.T, router http.Handler, email string) string {
	t.Helper()
	rec := performJSON(router, http.MethodPost, "/auth/register", gin.H{
		"name":     "Test User",
		"email":    email,
		"password": "secret123",
	}, "")
	if rec.Code != http.StatusCreated {
		t.Fatalf("register status = %d, body = %s", rec.Code, rec.Body.String())
	}

	var response struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode register response: %v", err)
	}
	if response.Token == "" {
		t.Fatal("expected token in register response")
	}
	return response.Token
}

func TestRegisterCreatesUserAndReturnsToken(t *testing.T) {
	router, db := newTestRouter(t)

	rec := performJSON(router, http.MethodPost, "/auth/register", gin.H{
		"name":     "Aida Student",
		"email":    "aida@example.com",
		"password": "secret123",
	}, "")

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}

	var response struct {
		Token string      `json:"token"`
		User  models.User `json:"user"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Token == "" || response.User.Email != "aida@example.com" {
		t.Fatalf("unexpected response: %+v", response)
	}

	var count int64
	db.Model(&models.User{}).Where("email = ?", "aida@example.com").Count(&count)
	if count != 1 {
		t.Fatalf("expected user persisted, count = %d", count)
	}
}

func TestRegisterRejectsDuplicateEmail(t *testing.T) {
	router, _ := newTestRouter(t)
	body := gin.H{"name": "Aida", "email": "aida@example.com", "password": "secret123"}

	first := performJSON(router, http.MethodPost, "/auth/register", body, "")
	if first.Code != http.StatusCreated {
		t.Fatalf("first register status = %d", first.Code)
	}
	second := performJSON(router, http.MethodPost, "/auth/register", body, "")
	if second.Code != http.StatusConflict {
		t.Fatalf("second register status = %d, body = %s", second.Code, second.Body.String())
	}
}

func TestLoginRejectsWrongPassword(t *testing.T) {
	router, _ := newTestRouter(t)
	registerAndToken(t, router, "login@example.com")

	rec := performJSON(router, http.MethodPost, "/auth/login", gin.H{
		"email":    "login@example.com",
		"password": "wrong-password",
	}, "")

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
}

func TestProtectedUsersRequiresJWT(t *testing.T) {
	router, _ := newTestRouter(t)

	rec := performJSON(router, http.MethodGet, "/users", nil, "")

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
}

func TestCreateAndListUsersWithJWT(t *testing.T) {
	router, _ := newTestRouter(t)
	token := registerAndToken(t, router, "admin@example.com")

	create := performJSON(router, http.MethodPost, "/users", gin.H{
		"name":      "Second User",
		"email":     "second@example.com",
		"is_active": true,
	}, token)
	if create.Code != http.StatusCreated {
		t.Fatalf("create status = %d, body = %s", create.Code, create.Body.String())
	}

	list := performJSON(router, http.MethodGet, "/users", nil, token)
	if list.Code != http.StatusOK {
		t.Fatalf("list status = %d, body = %s", list.Code, list.Body.String())
	}

	var users []models.User
	if err := json.Unmarshal(list.Body.Bytes(), &users); err != nil {
		t.Fatalf("decode users: %v", err)
	}
	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(users))
	}
}

func TestDeleteUserWithJWT(t *testing.T) {
	router, _ := newTestRouter(t)
	token := registerAndToken(t, router, "delete@example.com")

	rec := performJSON(router, http.MethodDelete, "/users/1", nil, token)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("delete status = %d, body = %s", rec.Code, rec.Body.String())
	}

	get := performJSON(router, http.MethodGet, "/users/1", nil, token)
	if get.Code != http.StatusNotFound {
		t.Fatalf("get deleted status = %d, body = %s", get.Code, get.Body.String())
	}
}
