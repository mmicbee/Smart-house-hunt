package properties

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"smart-house-hunt/internal/auth"
	"smart-house-hunt/middleware"
)

const testSecret = "test-secret-key"

func init() {
	gin.SetMode(gin.TestMode)
}

func newPropertiesTestRouter(repo *Repository, secret string) *gin.Engine {
	r := gin.New()
	handler := NewHandler(repo)

	r.GET("/properties", handler.ListApproved)
	r.GET("/properties/:id", handler.GetApproved)

	landlord := r.Group("/properties",
		middleware.RequireAuth(secret),
		middleware.RequireRole("landlord"),
	)
	{
		landlord.POST("", handler.Create)
		landlord.GET("/mine", handler.Mine)
		landlord.PUT("/:id", handler.Update)
		landlord.DELETE("/:id", handler.Delete)
		landlord.POST("/:id/images", handler.UploadImage)
	}

	return r
}

func generateTestToken(t *testing.T, userID int64, role string) string {
	t.Helper()
	token, err := auth.GenerateToken(testSecret, userID, role, time.Hour)
	if err != nil {
		t.Fatalf("failed to generate test token: %v", err)
	}
	return token
}

// 1. Ownership Checks (403 for another landlord)
func TestOwnershipChecks(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewRepository(db)
	router := newPropertiesTestRouter(repo, testSecret)

	// Create property owned by landlord 1
	prop := &Property{
		LandlordID:  1,
		Title:       "Landlord 1 Property",
		Description: "Nice flat",
		Location:    "Nairobi",
		Price:       25000,
		Bedrooms:    2,
		Bathrooms:   1,
		Status:      "pending",
	}
	if err := repo.Create(prop); err != nil {
		t.Fatalf("failed to seed property: %v", err)
	}

	landlord2Token := generateTestToken(t, 2, "landlord")

	// Another landlord attempts to update
	updateBody, _ := json.Marshal(PropertyRequest{
		Title:       "Hacked Title",
		Description: "Updated desc",
		Location:    "Nairobi",
		Price:       30000,
		Bedrooms:    2,
		Bathrooms:   1,
	})
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/properties/%d", prop.ID), bytes.NewReader(updateBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+landlord2Token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403 on update by another landlord, got %d (%s)", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "you can only edit your own listings") {
		t.Errorf("expected ownership error message, got %s", w.Body.String())
	}

	// Another landlord attempts to delete
	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/properties/%d", prop.ID), nil)
	req.Header.Set("Authorization", "Bearer "+landlord2Token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403 on delete by another landlord, got %d (%s)", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "you can only delete your own listings") {
		t.Errorf("expected ownership error message, got %s", w.Body.String())
	}

	// Another landlord attempts to upload image
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("image", "test.png")
	part.Write([]byte("image-data"))
	writer.Close()

	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/properties/%d/images", prop.ID), body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+landlord2Token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403 on upload by another landlord, got %d (%s)", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "you can only add images to your own listings") {
		t.Errorf("expected ownership error message, got %s", w.Body.String())
	}
}

// 2. 404 Not Found Checks
func TestNotFoundChecks(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewRepository(db)
	router := newPropertiesTestRouter(repo, testSecret)

	landlordToken := generateTestToken(t, 1, "landlord")

	// 404 on missing property for Update
	updateBody, _ := json.Marshal(PropertyRequest{
		Title:       "Valid Title",
		Description: "Valid desc",
		Location:    "Nairobi",
		Price:       20000,
		Bedrooms:    1,
		Bathrooms:   1,
	})
	req := httptest.NewRequest(http.MethodPut, "/properties/9999", bytes.NewReader(updateBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+landlordToken)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404 on update missing property, got %d", w.Code)
	}

	// 404 on missing property for Delete
	req = httptest.NewRequest(http.MethodDelete, "/properties/9999", nil)
	req.Header.Set("Authorization", "Bearer "+landlordToken)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404 on delete missing property, got %d", w.Code)
	}

	// 404 on missing property for UploadImage
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("image", "test.png")
	part.Write([]byte("img"))
	writer.Close()
	req = httptest.NewRequest(http.MethodPost, "/properties/9999/images", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+landlordToken)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404 on image upload missing property, got %d", w.Code)
	}

	// 404 on missing property for public Get
	req = httptest.NewRequest(http.MethodGet, "/properties/9999", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404 on public get missing property, got %d", w.Code)
	}

	// 404 on non-approved (pending) property for public Get
	pendingProp := &Property{
		LandlordID:  1,
		Title:       "Pending Prop",
		Description: "Pending desc",
		Location:    "Nairobi",
		Price:       20000,
		Bedrooms:    1,
		Bathrooms:   1,
		Status:      "pending",
	}
	_ = repo.Create(pendingProp)

	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/properties/%d", pendingProp.ID), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404 for unapproved property on GET /properties/:id, got %d", w.Code)
	}

	// 404 on rejected property
	rejectedProp := &Property{
		LandlordID:  1,
		Title:       "Rejected Prop",
		Description: "Rejected desc",
		Location:    "Nairobi",
		Price:       20000,
		Bedrooms:    1,
		Bathrooms:   1,
		Status:      "rejected",
	}
	_ = repo.Create(rejectedProp)

	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/properties/%d", rejectedProp.ID), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404 for rejected property, got %d", w.Code)
	}
}

// 3. 401 Unauthorized & Role Checks
func TestAuthAndRoleChecks(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewRepository(db)
	router := newPropertiesTestRouter(repo, testSecret)

	endpoints := []struct {
		method string
		path   string
	}{
		{http.MethodPost, "/properties"},
		{http.MethodGet, "/properties/mine"},
		{http.MethodPut, "/properties/1"},
		{http.MethodDelete, "/properties/1"},
		{http.MethodPost, "/properties/1/images"},
	}

	// Without token -> 401 Unauthorized
	for _, ep := range endpoints {
		req := httptest.NewRequest(ep.method, ep.path, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected 401 for anonymous on %s %s, got %d", ep.method, ep.path, w.Code)
		}
	}

	// With tenant role (wrong role) -> 403 Forbidden
	tenantToken := generateTestToken(t, 99, "tenant")
	for _, ep := range endpoints {
		req := httptest.NewRequest(ep.method, ep.path, nil)
		req.Header.Set("Authorization", "Bearer "+tenantToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusForbidden {
			t.Errorf("expected 403 for tenant on landlord route %s %s, got %d", ep.method, ep.path, w.Code)
		}
	}
}

// 4. Validation Errors (400 Bad Request)
func TestValidationErrors(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewRepository(db)
	router := newPropertiesTestRouter(repo, testSecret)

	landlordToken := generateTestToken(t, 1, "landlord")

	testCases := []struct {
		name        string
		req         PropertyRequest
		expectedErr string
	}{
		{
			name:        "empty title",
			req:         PropertyRequest{Title: "", Description: "Desc", Location: "Loc", Price: 10000},
			expectedErr: "title is required",
		},
		{
			name:        "empty description",
			req:         PropertyRequest{Title: "Title", Description: "", Location: "Loc", Price: 10000},
			expectedErr: "description is required",
		},
		{
			name:        "empty location",
			req:         PropertyRequest{Title: "Title", Description: "Desc", Location: "", Price: 10000},
			expectedErr: "location is required",
		},
		{
			name:        "zero price",
			req:         PropertyRequest{Title: "Title", Description: "Desc", Location: "Loc", Price: 0},
			expectedErr: "price must be greater than zero",
		},
		{
			name:        "negative price",
			req:         PropertyRequest{Title: "Title", Description: "Desc", Location: "Loc", Price: -500},
			expectedErr: "price must be greater than zero",
		},
		{
			name:        "negative bedrooms",
			req:         PropertyRequest{Title: "Title", Description: "Desc", Location: "Loc", Price: 10000, Bedrooms: -1},
			expectedErr: "bedrooms cannot be negative",
		},
		{
			name:        "negative bathrooms",
			req:         PropertyRequest{Title: "Title", Description: "Desc", Location: "Loc", Price: 10000, Bathrooms: -1},
			expectedErr: "bathrooms cannot be negative",
		},
	}

	for _, tc := range testCases {
		t.Run("Create "+tc.name, func(t *testing.T) {
			b, _ := json.Marshal(tc.req)
			req := httptest.NewRequest(http.MethodPost, "/properties", bytes.NewReader(b))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+landlordToken)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			if w.Code != http.StatusBadRequest {
				t.Errorf("expected 400, got %d (%s)", w.Code, w.Body.String())
			}
			if !strings.Contains(w.Body.String(), tc.expectedErr) {
				t.Errorf("expected error containing %q, got %s", tc.expectedErr, w.Body.String())
			}
		})
	}

	// Malformed JSON
	t.Run("Malformed JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/properties", strings.NewReader("not-json"))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+landlordToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400 for malformed json, got %d", w.Code)
		}
	})

	// Non-numeric ID parameters
	t.Run("Non-numeric IDs", func(t *testing.T) {
		// PUT
		req := httptest.NewRequest(http.MethodPut, "/properties/abc", strings.NewReader("{}"))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+landlordToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400 for non-numeric ID on PUT, got %d", w.Code)
		}

		// DELETE
		req = httptest.NewRequest(http.MethodDelete, "/properties/abc", nil)
		req.Header.Set("Authorization", "Bearer "+landlordToken)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400 for non-numeric ID on DELETE, got %d", w.Code)
		}

		// GET
		req = httptest.NewRequest(http.MethodGet, "/properties/abc", nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400 for non-numeric ID on GET, got %d", w.Code)
		}

		// POST images
		req = httptest.NewRequest(http.MethodPost, "/properties/abc/images", nil)
		req.Header.Set("Authorization", "Bearer "+landlordToken)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400 for non-numeric ID on image upload, got %d", w.Code)
		}
	})

	// Image upload validations
	t.Run("Image upload invalid extension", func(t *testing.T) {
		prop := &Property{
			LandlordID:  1,
			Title:       "Photo Prop",
			Description: "Desc",
			Location:    "Loc",
			Price:       10000,
		}
		_ = repo.Create(prop)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("image", "document.pdf")
		part.Write([]byte("%PDF-fake"))
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/properties/%d/images", prop.ID), body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("Authorization", "Bearer "+landlordToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400 for pdf image, got %d", w.Code)
		}
		if !strings.Contains(w.Body.String(), "only jpg, jpeg, png, webp, or gif images are supported") {
			t.Errorf("expected unsupported extension error, got %s", w.Body.String())
		}
	})

	t.Run("Image upload missing file", func(t *testing.T) {
		prop := &Property{
			LandlordID:  1,
			Title:       "Photo Prop 2",
			Description: "Desc",
			Location:    "Loc",
			Price:       10000,
		}
		_ = repo.Create(prop)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/properties/%d/images", prop.ID), body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("Authorization", "Bearer "+landlordToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400 for missing file, got %d", w.Code)
		}
	})
}

// 5. Full lifecycle and approved visibility tests
func TestPropertyLifecycleAndApprovedListing(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewRepository(db)
	router := newPropertiesTestRouter(repo, testSecret)

	defer os.RemoveAll("storage")

	landlordToken := generateTestToken(t, 1, "landlord")

	// 1. Create property
	createJSON := `{
		"title": "Modern 2-Bedroom Apartment",
		"description": "Quiet neighborhood with 24/7 security",
		"location": "Kakamega CBD",
		"price": 22000,
		"bedrooms": 2,
		"bathrooms": 1,
		"amenities": "wifi,parking,water"
	}`
	req := httptest.NewRequest(http.MethodPost, "/properties", strings.NewReader(createJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+landlordToken)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 on create, got %d (%s)", w.Code, w.Body.String())
	}
	var created PropertyResponse
	if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil {
		t.Fatalf("failed to decode created property: %v", err)
	}
	if created.ID == 0 || created.Status != "pending" {
		t.Errorf("unexpected created property response: %+v", created)
	}

	// 2. View own listings (GET /properties/mine)
	req = httptest.NewRequest(http.MethodGet, "/properties/mine", nil)
	req.Header.Set("Authorization", "Bearer "+landlordToken)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on /properties/mine, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "Modern 2-Bedroom Apartment") {
		t.Errorf("expected created property in /properties/mine, got %s", w.Body.String())
	}

	// 3. Public list (GET /properties) should NOT show pending property
	req = httptest.NewRequest(http.MethodGet, "/properties", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on GET /properties, got %d", w.Code)
	}
	if strings.Contains(w.Body.String(), "Modern 2-Bedroom Apartment") {
		t.Errorf("pending property should not appear in public list: %s", w.Body.String())
	}

	// 4. Public details (GET /properties/:id) should return 404 for pending
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/properties/%d", created.ID), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404 for pending property details, got %d", w.Code)
	}

	// 5. Update property status directly to 'approved' (as an admin would)
	_, err := db.Exec("UPDATE properties SET status = 'approved' WHERE id = ?", created.ID)
	if err != nil {
		t.Fatalf("failed to approve property: %v", err)
	}

	// 6. Upload an image to the approved property
	imgBody := &bytes.Buffer{}
	imgWriter := multipart.NewWriter(imgBody)
	part, _ := imgWriter.CreateFormFile("image", "house.jpg")
	part.Write([]byte("fake-jpeg-data"))
	imgWriter.Close()

	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/properties/%d/images", created.ID), imgBody)
	req.Header.Set("Content-Type", imgWriter.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+landlordToken)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 on image upload, got %d (%s)", w.Code, w.Body.String())
	}

	// 7. Public list (GET /properties) should now show the approved property
	req = httptest.NewRequest(http.MethodGet, "/properties", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on GET /properties, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "Modern 2-Bedroom Apartment") {
		t.Errorf("approved property missing from public list: %s", w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "house.jpg") && !strings.Contains(w.Body.String(), "uploads/property-images") {
		t.Errorf("uploaded image missing from public property item: %s", w.Body.String())
	}

	// 8. Public details (GET /properties/:id) should return 200 with property and images
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/properties/%d", created.ID), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on approved property details, got %d (%s)", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "Modern 2-Bedroom Apartment") {
		t.Errorf("property title missing from details: %s", w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "uploads/property-images") {
		t.Errorf("property image missing from details: %s", w.Body.String())
	}

	// 9. Update property
	updateJSON := `{
		"title": "Renovated 2-Bedroom Apartment",
		"description": "Quiet neighborhood with 24/7 security",
		"location": "Kakamega CBD",
		"price": 24000,
		"bedrooms": 2,
		"bathrooms": 2,
		"amenities": "wifi,parking,water,balcony"
	}`
	req = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/properties/%d", created.ID), strings.NewReader(updateJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+landlordToken)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on update, got %d (%s)", w.Code, w.Body.String())
	}

	// 10. Delete property
	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/properties/%d", created.ID), nil)
	req.Header.Set("Authorization", "Bearer "+landlordToken)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on delete, got %d (%s)", w.Code, w.Body.String())
	}

	// 11. Public details after delete returns 404
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/properties/%d", created.ID), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404 after property deletion, got %d", w.Code)
	}
}
