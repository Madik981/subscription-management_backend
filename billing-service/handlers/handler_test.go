package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"

	"subscription-management/billing-service/models"
)

const testJWTSecret = "test-secret"

func newTestRouter(t *testing.T, accountsURL string) (*gin.Engine, *gorm.DB) {
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
	if err := db.AutoMigrate(&models.Plan{}, &models.Billing{}); err != nil {
		t.Fatalf("migrate test db: %v", err)
	}

	router := gin.New()
	NewHandler(db, testJWTSecret, accountsURL, "internal-token").RegisterRoutes(router)
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

func testToken(t *testing.T) string {
	t.Helper()
	now := time.Now().UTC()
	claims := AuthClaims{
		UserID: 1,
		Email:  "tester@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "user",
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(testJWTSecret))
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	return token
}

func mockAccountsService(t *testing.T, ok bool) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/internal/users/42" {
			http.NotFound(w, r)
			return
		}
		if r.Header.Get("X-Internal-Token") != "internal-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":42}`))
	}))
}

func TestCreatePlanDefaultsCurrencyAndCycle(t *testing.T) {
	router, _ := newTestRouter(t, "")
	token := testToken(t)

	rec := performJSON(router, http.MethodPost, "/plans", gin.H{
		"name":  "Starter",
		"price": 9.99,
	}, token)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}

	var plan models.Plan
	if err := json.Unmarshal(rec.Body.Bytes(), &plan); err != nil {
		t.Fatalf("decode plan: %v", err)
	}
	if plan.Currency != "USD" || plan.BillingCycle != "monthly" {
		t.Fatalf("unexpected defaults: %+v", plan)
	}
}

func TestListPlansRequiresJWT(t *testing.T) {
	router, _ := newTestRouter(t, "")

	rec := performJSON(router, http.MethodGet, "/plans", nil, "")

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
}

func TestUpdatePlanWithJWT(t *testing.T) {
	router, db := newTestRouter(t, "")
	token := testToken(t)
	plan := models.Plan{Name: "Pro", Price: 19, Currency: "USD", BillingCycle: "monthly"}
	if err := db.Create(&plan).Error; err != nil {
		t.Fatalf("seed plan: %v", err)
	}

	rec := performJSON(router, http.MethodPatch, "/plans/1", gin.H{
		"price":         29.99,
		"billing_cycle": "yearly",
	}, token)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}

	var updated models.Plan
	if err := json.Unmarshal(rec.Body.Bytes(), &updated); err != nil {
		t.Fatalf("decode plan: %v", err)
	}
	if updated.Price != 29.99 || updated.BillingCycle != "yearly" {
		t.Fatalf("unexpected updated plan: %+v", updated)
	}
}

func TestCreateBillingUsesPlanPriceAndRemoteUser(t *testing.T) {
	accounts := mockAccountsService(t, true)
	defer accounts.Close()

	router, db := newTestRouter(t, accounts.URL)
	token := testToken(t)
	plan := models.Plan{Name: "Pro", Price: 25, Currency: "USD", BillingCycle: "monthly"}
	if err := db.Create(&plan).Error; err != nil {
		t.Fatalf("seed plan: %v", err)
	}

	rec := performJSON(router, http.MethodPost, "/billings", gin.H{
		"user_id":     42,
		"plan_id":     plan.ID,
		"due_date":    "2026-05-10T00:00:00Z",
		"description": "May invoice",
	}, token)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}

	var billing models.Billing
	if err := json.Unmarshal(rec.Body.Bytes(), &billing); err != nil {
		t.Fatalf("decode billing: %v", err)
	}
	if billing.Amount != plan.Price || billing.Status != models.BillingStatusPending {
		t.Fatalf("unexpected billing: %+v", billing)
	}
}

func TestPayBillingMarksPaid(t *testing.T) {
	router, db := newTestRouter(t, "")
	token := testToken(t)
	plan := models.Plan{Name: "Pro", Price: 25, Currency: "USD", BillingCycle: "monthly"}
	if err := db.Create(&plan).Error; err != nil {
		t.Fatalf("seed plan: %v", err)
	}
	billing := models.Billing{
		UserID:  42,
		PlanID:  plan.ID,
		Amount:  25,
		Status:  models.BillingStatusPending,
		DueDate: time.Now().UTC().Add(24 * time.Hour),
	}
	if err := db.Create(&billing).Error; err != nil {
		t.Fatalf("seed billing: %v", err)
	}

	rec := performJSON(router, http.MethodPatch, "/billings/1/pay", nil, token)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}

	var paid models.Billing
	if err := json.Unmarshal(rec.Body.Bytes(), &paid); err != nil {
		t.Fatalf("decode billing: %v", err)
	}
	if paid.Status != models.BillingStatusPaid || paid.PaidAt == nil {
		t.Fatalf("expected paid billing, got %+v", paid)
	}
}

func TestFailAndDeleteBilling(t *testing.T) {
	router, db := newTestRouter(t, "")
	token := testToken(t)
	plan := models.Plan{Name: "Pro", Price: 25, Currency: "USD", BillingCycle: "monthly"}
	if err := db.Create(&plan).Error; err != nil {
		t.Fatalf("seed plan: %v", err)
	}
	billing := models.Billing{
		UserID:  42,
		PlanID:  plan.ID,
		Amount:  25,
		Status:  models.BillingStatusPending,
		DueDate: time.Now().UTC().Add(24 * time.Hour),
	}
	if err := db.Create(&billing).Error; err != nil {
		t.Fatalf("seed billing: %v", err)
	}

	fail := performJSON(router, http.MethodPatch, "/billings/1/fail", nil, token)
	if fail.Code != http.StatusOK {
		t.Fatalf("fail status = %d, body = %s", fail.Code, fail.Body.String())
	}

	var failed models.Billing
	if err := json.Unmarshal(fail.Body.Bytes(), &failed); err != nil {
		t.Fatalf("decode failed billing: %v", err)
	}
	if failed.Status != models.BillingStatusFailed {
		t.Fatalf("expected failed status, got %+v", failed)
	}

	deleted := performJSON(router, http.MethodDelete, "/billings/1", nil, token)
	if deleted.Code != http.StatusNoContent {
		t.Fatalf("delete status = %d, body = %s", deleted.Code, deleted.Body.String())
	}
}
