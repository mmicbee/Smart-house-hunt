package favorites

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"smart-house-hunt/internal/auth"
	"smart-house-hunt/internal/properties"
	"smart-house-hunt/middleware"
)

const testSecret = "test-secret-key"

func init() {
	gin.SetMode(gin.TestMode)
}

func newFavoritesTestRouter(repo *Repository, propRepo *properties.Repository) *gin.Engine {
	r := gin.New()
	handler := NewHandler(repo, propRepo)

	favs := r.Group("/favorites", middleware.RequireAuth(testSecret))
	{
		favs.GET("", handler.List)
		favs.GET("/ids", handler.FavoriteIDs)
		favs.GET("/check/:id", handler.Check)
		favs.POST("", handler.Add)
		favs.POST("/:id", handler.Add)
		favs.DELETE("/:id", handler.Remove)
	}

	return r
}

func generateToken(t *testing.T, userID int64, role string) string {
	t.Helper()
	token, err := auth.GenerateToken(testSecret, userID, role, time.Hour)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}
	return token
}

func TestFavoritesHandlers(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	favRepo := NewRepository(db)
	propRepo := properties.NewRepository(db)
	router := newFavoritesTestRouter(favRepo, propRepo)

	p1 := &properties.Property{LandlordID: 1, Title: "Kileleshwa Haven", Description: "Spacious", Location: "Kileleshwa", Price: 45000, Bedrooms: 3, Bathrooms: 2, Amenities: "wifi, parking", Status: "approved"}
	p2 := &properties.Property{LandlordID: 1, Title: "Runda Villa", Description: "Luxury", Location: "Runda", Price: 120000, Bedrooms: 5, Bathrooms: 4, Amenities: "pool, garden", Status: "approved"}
	pPending := &properties.Property{LandlordID: 1, Title: "Unapproved House", Description: "Waiting", Location: "CBD", Price: 20000, Status: "pending"}
	_ = propRepo.Create(p1)
	_ = propRepo.Create(p2)
	_ = propRepo.Create(pPending)

	// Add an image to p1
	_ = propRepo.AddImage(&properties.PropertyImage{PropertyID: p1.ID, FilePath: "uploads/property-images/1/test.jpg"})

	tenantID := int64(42)
	tenantToken := generateToken(t, tenantID, "tenant")

	// 1. Unauthenticated request returns 401
	unauthReq := httptest.NewRequest(http.MethodGet, "/favorites", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, unauthReq)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 for unauthenticated request, got %d", w.Code)
	}

	// 2. Add non-existent property returns 404
	body, _ := json.Marshal(map[string]int64{"property_id": 9999})
	req := httptest.NewRequest(http.MethodPost, "/favorites", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tenantToken)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404 for non-existent property, got %d", w.Code)
	}

	// 3. Add pending (unapproved) property returns 404
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/favorites/%d", pPending.ID), nil)
	req.Header.Set("Authorization", "Bearer "+tenantToken)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404 for unapproved property, got %d", w.Code)
	}

	// 4. Save approved property p1 (via JSON body)
	body, _ = json.Marshal(map[string]int64{"property_id": p1.ID})
	req = httptest.NewRequest(http.MethodPost, "/favorites", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tenantToken)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 on saving favorite, got %d (%s)", w.Code, w.Body.String())
	}

	// 5. Save again (duplicate) returns 200 with notice
	req = httptest.NewRequest(http.MethodPost, "/favorites", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tenantToken)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on duplicate favorite, got %d", w.Code)
	}

	// 6. Save p2 via URL param POST /favorites/:id
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/favorites/%d", p2.ID), nil)
	req.Header.Set("Authorization", "Bearer "+tenantToken)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 on saving favorite via URL, got %d", w.Code)
	}

	// 7. Check favorite status for p1 and pPending
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/favorites/check/%d", p1.ID), nil)
	req.Header.Set("Authorization", "Bearer "+tenantToken)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for check, got %d", w.Code)
	}
	var checkResp struct {
		IsFavorite bool `json:"is_favorite"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &checkResp)
	if !checkResp.IsFavorite {
		t.Errorf("expected is_favorite = true for p1")
	}

	// 8. List favorite IDs
	req = httptest.NewRequest(http.MethodGet, "/favorites/ids", nil)
	req.Header.Set("Authorization", "Bearer "+tenantToken)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for ids, got %d", w.Code)
	}
	var idsResp struct {
		FavoriteIDs []int64 `json:"favorite_ids"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &idsResp)
	if len(idsResp.FavoriteIDs) != 2 {
		t.Errorf("expected 2 favorite ids, got %d", len(idsResp.FavoriteIDs))
	}

	// 9. List favorite properties (with images)
	req = httptest.NewRequest(http.MethodGet, "/favorites", nil)
	req.Header.Set("Authorization", "Bearer "+tenantToken)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for list, got %d", w.Code)
	}
	var listResp struct {
		Properties []properties.PropertyListItem `json:"properties"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &listResp)
	if len(listResp.Properties) != 2 {
		t.Fatalf("expected 2 favorite properties, got %d", len(listResp.Properties))
	}
	// Verify p1 has image attached
	foundImage := false
	for _, item := range listResp.Properties {
		if item.ID == p1.ID && len(item.Images) > 0 {
			foundImage = true
		}
	}
	if !foundImage {
		t.Errorf("expected p1 to include its property images")
	}

	// 10. Remove favorite p1
	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/favorites/%d", p1.ID), nil)
	req.Header.Set("Authorization", "Bearer "+tenantToken)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for remove, got %d", w.Code)
	}

	// 11. Remove p1 again returns 404
	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/favorites/%d", p1.ID), nil)
	req.Header.Set("Authorization", "Bearer "+tenantToken)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 on second remove, got %d", w.Code)
	}
}
