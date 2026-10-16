package admin

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"smart-house-hunt/internal/auth"
	"smart-house-hunt/middleware"
)

const adminTestSecret = "test-admin-secret-123456789012345"

func setupTestRouter(t *testing.T) (*gin.Engine, *Repository) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	db := setupTestDB(t)

	_, err := db.Exec(`
		INSERT INTO users (id, role_id, name, email, password, status) VALUES
		(1, 4, 'Admin Bob', 'admin@test.com', 'pass', 'active'),
		(2, 2, 'Landlord Larry', 'larry@test.com', 'pass', 'active'),
		(3, 1, 'Tenant Tim', 'tim@test.com', 'pass', 'active');

		INSERT INTO properties (id, landlord_id, title, description, location, price, bedrooms, bathrooms, amenities, property_type, status) VALUES
		(1, 2, 'Pending Villa', 'Big house', 'Karen', 120000, 4, 3, 'Pool, Garden', 'house', 'pending'),
		(2, 2, 'Approved Flat', 'Modern flat', 'Kilimani', 40000, 2, 2, 'WiFi', 'apartment', 'approved');
	`)
	if err != nil {
		t.Fatalf("failed to insert initial test data: %v", err)
	}

	repo := NewRepository(db)
	handler := NewHandler(repo)

	r := gin.New()
	adminAPI := r.Group("/api/admin",
		middleware.RequireAuth(adminTestSecret),
		middleware.RequireRole("admin"),
	)
	{
		adminAPI.GET("/reports", handler.GetReports)
		adminAPI.GET("/users", handler.ListUsers)
		adminAPI.PUT("/users/:id/status", handler.UpdateUserStatus)

		adminAPI.GET("/properties", handler.ListProperties)
		adminAPI.GET("/properties/:id", handler.GetProperty)
		adminAPI.PUT("/properties/:id/status", handler.UpdatePropertyStatus)
		adminAPI.PUT("/properties/:id/approve", handler.ApproveProperty)
		adminAPI.PUT("/properties/:id/reject", handler.RejectProperty)
		adminAPI.DELETE("/properties/:id", handler.DeleteProperty)
	}

	return r, repo
}

func adminAuthToken(t *testing.T, userID int64, role string) string {
	t.Helper()
	token, err := auth.GenerateToken(adminTestSecret, userID, role, time.Hour)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}
	return "Bearer " + token
}

func TestAdminListUsersEndpoint(t *testing.T) {
	r, _ := setupTestRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/api/admin/users", nil)
	req.Header.Set("Authorization", adminAuthToken(t, 1, "admin"))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var res struct {
		Users []AdminUser `json:"users"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to unmarshal json: %v", err)
	}
	if len(res.Users) != 3 {
		t.Fatalf("expected 3 users, got %d", len(res.Users))
	}
}

func TestAdminUpdateUserStatusEndpoint(t *testing.T) {
	r, repo := setupTestRouter(t)

	// 1. Non-admin should be forbidden (403)
	body, _ := json.Marshal(UpdateUserStatusRequest{Status: "suspended"})
	req := httptest.NewRequest(http.MethodPut, "/api/admin/users/3/status", bytes.NewReader(body))
	req.Header.Set("Authorization", adminAuthToken(t, 3, "tenant"))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for tenant, got %d", w.Code)
	}

	// 2. Admin cannot suspend self (400)
	req = httptest.NewRequest(http.MethodPut, "/api/admin/users/1/status", bytes.NewReader(body))
	req.Header.Set("Authorization", adminAuthToken(t, 1, "admin"))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 when suspending self, got %d: %s", w.Code, w.Body.String())
	}

	// 3. Admin suspends tenant (200)
	req = httptest.NewRequest(http.MethodPut, "/api/admin/users/3/status", bytes.NewReader(body))
	req.Header.Set("Authorization", adminAuthToken(t, 1, "admin"))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 when suspending tenant, got %d: %s", w.Code, w.Body.String())
	}

	u, _ := repo.GetUserByID(3)
	if u.Status != "suspended" {
		t.Fatalf("expected status 'suspended', got '%s'", u.Status)
	}

	// 4. Invalid status rejected (400)
	badBody, _ := json.Marshal(UpdateUserStatusRequest{Status: "unknown-status"})
	req = httptest.NewRequest(http.MethodPut, "/api/admin/users/3/status", bytes.NewReader(badBody))
	req.Header.Set("Authorization", adminAuthToken(t, 1, "admin"))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid status, got %d", w.Code)
	}
}

func TestAdminPropertyModeration(t *testing.T) {
	r, repo := setupTestRouter(t)

	// 1. List properties with pending filter
	req := httptest.NewRequest(http.MethodGet, "/api/admin/properties?status=pending", nil)
	req.Header.Set("Authorization", adminAuthToken(t, 1, "admin"))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var listRes struct {
		Properties []AdminProperty `json:"properties"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &listRes)
	if len(listRes.Properties) != 1 || listRes.Properties[0].ID != 1 {
		t.Fatalf("expected 1 pending property, got %d", len(listRes.Properties))
	}

	// 2. Approve property 1
	req = httptest.NewRequest(http.MethodPut, "/api/admin/properties/1/approve", nil)
	req.Header.Set("Authorization", adminAuthToken(t, 1, "admin"))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on approve, got %d: %s", w.Code, w.Body.String())
	}

	p1, _ := repo.GetPropertyByID(1)
	if p1.Status != "approved" {
		t.Fatalf("expected status 'approved', got '%s'", p1.Status)
	}

	// 3. Reject property 1
	req = httptest.NewRequest(http.MethodPut, "/api/admin/properties/1/reject", nil)
	req.Header.Set("Authorization", adminAuthToken(t, 1, "admin"))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on reject, got %d: %s", w.Code, w.Body.String())
	}

	p1AfterReject, _ := repo.GetPropertyByID(1)
	if p1AfterReject.Status != "rejected" {
		t.Fatalf("expected status 'rejected', got '%s'", p1AfterReject.Status)
	}

	// 4. Update status with body
	body, _ := json.Marshal(UpdatePropertyStatusRequest{Status: "suspended"})
	req = httptest.NewRequest(http.MethodPut, "/api/admin/properties/1/status", bytes.NewReader(body))
	req.Header.Set("Authorization", adminAuthToken(t, 1, "admin"))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on status update, got %d", w.Code)
	}

	// 5. Delete property
	req = httptest.NewRequest(http.MethodDelete, "/api/admin/properties/1", nil)
	req.Header.Set("Authorization", adminAuthToken(t, 1, "admin"))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on delete, got %d", w.Code)
	}

	// Deleting again should return 404
	req = httptest.NewRequest(http.MethodDelete, "/api/admin/properties/1", nil)
	req.Header.Set("Authorization", adminAuthToken(t, 1, "admin"))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 on deleting non-existent property, got %d", w.Code)
	}
}

func TestAdminGetReportsEndpoint(t *testing.T) {
	r, _ := setupTestRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/api/admin/reports", nil)
	req.Header.Set("Authorization", adminAuthToken(t, 1, "admin"))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var res struct {
		Reports PlatformReports `json:"reports"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to unmarshal reports: %v", err)
	}
	if res.Reports.Users.Total != 3 {
		t.Errorf("expected 3 total users, got %d", res.Reports.Users.Total)
	}
	if res.Reports.Properties.Total != 2 {
		t.Errorf("expected 2 total properties, got %d", res.Reports.Properties.Total)
	}
}

func TestAdminUnauthorizedAccess(t *testing.T) {
	r, _ := setupTestRouter(t)

	// No token provided -> 401
	endpoints := []string{
		"/api/admin/reports",
		"/api/admin/users",
		"/api/admin/properties",
	}

	for _, ep := range endpoints {
		req := httptest.NewRequest(http.MethodGet, ep, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected 401 for unauthenticated request to %s, got %d", ep, w.Code)
		}
	}
}
