package viewing

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

func newViewingTestRouter(repo *Repository, propRepo *properties.Repository) *gin.Engine {
	r := gin.New()
	handler := NewHandler(repo, propRepo)

	viewings := r.Group("/viewings", middleware.RequireAuth(testSecret))
	{
		viewings.POST("", handler.Create)
		viewings.GET("/mine", handler.ListMine)
		viewings.GET("/landlord", handler.ListLandlord)
		viewings.GET("/:id", handler.GetByID)
		viewings.PATCH("/:id/status", handler.UpdateStatus)
		viewings.POST("/:id/accept", handler.Accept)
		viewings.POST("/:id/reject", handler.Reject)
		viewings.POST("/:id/cancel", handler.Cancel)
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

func TestViewingHandlers(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewRepository(db)
	propRepo := properties.NewRepository(db)
	router := newViewingTestRouter(repo, propRepo)

	// User 1 is landlord, User 2 is tenant Alice, User 3 is tenant Bob
	landlordToken := generateToken(t, 1, "landlord")
	aliceToken := generateToken(t, 2, "tenant")
	bobToken := generateToken(t, 3, "tenant")

	prop := &properties.Property{
		LandlordID:  1,
		Title:       "Westlands Studio",
		Description: "Cozy and convenient",
		Location:    "Westlands",
		Price:       35000,
		Bedrooms:    1,
		Bathrooms:   1,
		Amenities:   "wifi,security",
		Status:      "approved",
	}
	if err := propRepo.Create(prop); err != nil {
		t.Fatalf("failed to create property: %v", err)
	}

	unapprovedProp := &properties.Property{
		LandlordID:  1,
		Title:       "Pending Condo",
		Description: "Not yet verified",
		Location:    "Parklands",
		Price:       50000,
		Status:      "pending",
	}
	if err := propRepo.Create(unapprovedProp); err != nil {
		t.Fatalf("failed to create unapproved property: %v", err)
	}

	t.Run("Create viewing request validation", func(t *testing.T) {
		// Missing auth
		req := httptest.NewRequest(http.MethodPost, "/viewings", bytes.NewBufferString(`{}`))
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected 401, got %d", rec.Code)
		}

		// Missing property_id
		body, _ := json.Marshal(map[string]any{"preferred_date": "2026-10-20"})
		req = httptest.NewRequest(http.MethodPost, "/viewings", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+aliceToken)
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400 for missing property_id, got %d", rec.Code)
		}

		// Missing preferred_date
		body, _ = json.Marshal(map[string]any{"property_id": prop.ID})
		req = httptest.NewRequest(http.MethodPost, "/viewings", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+aliceToken)
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400 for missing preferred_date, got %d", rec.Code)
		}

		// Non-existent property
		body, _ = json.Marshal(map[string]any{"property_id": 9999, "preferred_date": "2026-10-20"})
		req = httptest.NewRequest(http.MethodPost, "/viewings", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+aliceToken)
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Errorf("expected 404 for non-existent property, got %d", rec.Code)
		}

		// Unapproved property
		body, _ = json.Marshal(map[string]any{"property_id": unapprovedProp.ID, "preferred_date": "2026-10-20"})
		req = httptest.NewRequest(http.MethodPost, "/viewings", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+aliceToken)
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Errorf("expected 404 for unapproved property, got %d", rec.Code)
		}

		// Landlord requesting viewing on own property
		body, _ = json.Marshal(map[string]any{"property_id": prop.ID, "preferred_date": "2026-10-20"})
		req = httptest.NewRequest(http.MethodPost, "/viewings", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+landlordToken)
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400 when landlord requests own property, got %d", rec.Code)
		}
	})

	var createdID int64
	t.Run("Create viewing request success and duplicate prevention", func(t *testing.T) {
		body, _ := json.Marshal(map[string]any{
			"property_id":    prop.ID,
			"preferred_date": "2026-10-22T10:00",
			"message":        "Looking to move next month",
		})
		req := httptest.NewRequest(http.MethodPost, "/viewings", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+aliceToken)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
		}

		var created ViewingRequest
		if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if created.ID == 0 || created.Status != StatusPending {
			t.Fatalf("unexpected created request: %+v", created)
		}
		createdID = created.ID

		// Duplicate request should return 409 Conflict
		reqDup := httptest.NewRequest(http.MethodPost, "/viewings", bytes.NewBuffer(body))
		reqDup.Header.Set("Authorization", "Bearer "+aliceToken)
		recDup := httptest.NewRecorder()
		router.ServeHTTP(recDup, reqDup)
		if recDup.Code != http.StatusConflict {
			t.Errorf("expected 409 for duplicate request, got %d", recDup.Code)
		}
	})

	t.Run("ListMine and ListLandlord", func(t *testing.T) {
		// Alice checks mine
		req := httptest.NewRequest(http.MethodGet, "/viewings/mine", nil)
		req.Header.Set("Authorization", "Bearer "+aliceToken)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}
		var mineResp struct {
			Viewings []ViewingRequestDetail `json:"viewing_requests"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &mineResp); err != nil {
			t.Fatalf("failed to parse mine: %v", err)
		}
		if len(mineResp.Viewings) != 1 {
			t.Fatalf("expected 1 viewing for Alice, got %d", len(mineResp.Viewings))
		}

		// Landlord checks landlord viewings
		reqL := httptest.NewRequest(http.MethodGet, "/viewings/landlord", nil)
		reqL.Header.Set("Authorization", "Bearer "+landlordToken)
		recL := httptest.NewRecorder()
		router.ServeHTTP(recL, reqL)
		if recL.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", recL.Code)
		}
		var landlordResp struct {
			Viewings []ViewingRequestDetail `json:"viewing_requests"`
		}
		if err := json.Unmarshal(recL.Body.Bytes(), &landlordResp); err != nil {
			t.Fatalf("failed to parse landlord: %v", err)
		}
		if len(landlordResp.Viewings) != 1 {
			t.Fatalf("expected 1 viewing for landlord, got %d", len(landlordResp.Viewings))
		}
	})

	t.Run("GetByID permissions", func(t *testing.T) {
		// Bob (uninvolved) gets 403
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/viewings/%d", createdID), nil)
		req.Header.Set("Authorization", "Bearer "+bobToken)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		if rec.Code != http.StatusForbidden {
			t.Errorf("expected 403 for uninvolved user, got %d", rec.Code)
		}

		// Alice (tenant) gets 200
		req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/viewings/%d", createdID), nil)
		req.Header.Set("Authorization", "Bearer "+aliceToken)
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Errorf("expected 200 for tenant, got %d", rec.Code)
		}

		// Landlord gets 200
		req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/viewings/%d", createdID), nil)
		req.Header.Set("Authorization", "Bearer "+landlordToken)
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Errorf("expected 200 for landlord, got %d", rec.Code)
		}
	})

	t.Run("Update status authorization and flow", func(t *testing.T) {
		// Alice tries to accept (forbidden, only landlord can accept)
		body, _ := json.Marshal(map[string]any{"status": "accepted"})
		req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/viewings/%d/status", createdID), bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+aliceToken)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		if rec.Code != http.StatusForbidden {
			t.Errorf("expected 403 when tenant tries to accept, got %d", rec.Code)
		}

		// Landlord accepts via accept endpoint
		reqAccept := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/viewings/%d/accept", createdID), nil)
		reqAccept.Header.Set("Authorization", "Bearer "+landlordToken)
		recAccept := httptest.NewRecorder()
		router.ServeHTTP(recAccept, reqAccept)
		if recAccept.Code != http.StatusOK {
			t.Fatalf("expected 200 on accept, got %d: %s", recAccept.Code, recAccept.Body.String())
		}

		var updated ViewingRequestDetail
		if err := json.Unmarshal(recAccept.Body.Bytes(), &updated); err != nil {
			t.Fatalf("failed to decode updated viewing: %v", err)
		}
		if updated.Status != StatusAccepted {
			t.Errorf("expected status accepted, got %s", updated.Status)
		}

		// Landlord can reject
		reqReject := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/viewings/%d/reject", createdID), nil)
		reqReject.Header.Set("Authorization", "Bearer "+landlordToken)
		recReject := httptest.NewRecorder()
		router.ServeHTTP(recReject, reqReject)
		if recReject.Code != http.StatusOK {
			t.Fatalf("expected 200 on reject, got %d", recReject.Code)
		}

		// Create a second request from Bob to test cancel
		bodyBob, _ := json.Marshal(map[string]any{
			"property_id":    prop.ID,
			"preferred_date": "2026-10-23T15:00",
		})
		reqBob := httptest.NewRequest(http.MethodPost, "/viewings", bytes.NewBuffer(bodyBob))
		reqBob.Header.Set("Authorization", "Bearer "+bobToken)
		recBob := httptest.NewRecorder()
		router.ServeHTTP(recBob, reqBob)
		if recBob.Code != http.StatusCreated {
			t.Fatalf("expected 201 for bob, got %d", recBob.Code)
		}
		var bobReq ViewingRequest
		_ = json.Unmarshal(recBob.Body.Bytes(), &bobReq)

		// Alice tries to cancel Bob's request (forbidden)
		reqCancel := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/viewings/%d/cancel", bobReq.ID), nil)
		reqCancel.Header.Set("Authorization", "Bearer "+aliceToken)
		recCancel := httptest.NewRecorder()
		router.ServeHTTP(recCancel, reqCancel)
		if recCancel.Code != http.StatusForbidden {
			t.Errorf("expected 403 when Alice tries to cancel Bob's request, got %d", recCancel.Code)
		}

		// Bob cancels his own request (success)
		reqCancelBob := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/viewings/%d/cancel", bobReq.ID), nil)
		reqCancelBob.Header.Set("Authorization", "Bearer "+bobToken)
		recCancelBob := httptest.NewRecorder()
		router.ServeHTTP(recCancelBob, reqCancelBob)
		if recCancelBob.Code != http.StatusOK {
			t.Fatalf("expected 200 when Bob cancels, got %d", recCancelBob.Code)
		}
	})
}
