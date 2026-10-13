package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"smart-house-hunt/internal/auth"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func newAuthTestRouter(secret string) *gin.Engine {
	r := gin.New()
	protected := r.Group("/protected", RequireAuth(secret))
	protected.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"user_id": c.GetInt64(ContextUserID),
			"role":    c.GetString(ContextRole),
		})
	})
	return r
}

func TestRequireAuthValidToken(t *testing.T) {
	secret := "test-secret"
	token, err := auth.GenerateToken(secret, 7, "tenant", time.Hour)
	if err != nil {
		t.Fatalf("GenerateToken returned an error: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	newAuthTestRouter(secret).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d (%s)", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), `"user_id":7`) {
		t.Errorf("expected user_id 7 in context, got %s", w.Body.String())
	}
	if !strings.Contains(w.Body.String(), `"role":"tenant"`) {
		t.Errorf("expected role tenant in context, got %s", w.Body.String())
	}
}

func TestRequireAuthMissingHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()

	newAuthTestRouter("test-secret").ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestRequireAuthMalformedHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Token abc123")
	w := httptest.NewRecorder()

	newAuthTestRouter("test-secret").ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestRequireAuthInvalidToken(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer not-a-valid-token")
	w := httptest.NewRecorder()

	newAuthTestRouter("test-secret").ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestRequireAuthExpiredToken(t *testing.T) {
	secret := "test-secret"
	token, err := auth.GenerateToken(secret, 7, "tenant", -time.Minute)
	if err != nil {
		t.Fatalf("GenerateToken returned an error: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	newAuthTestRouter(secret).ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestRequireAuthWrongSecret(t *testing.T) {
	token, err := auth.GenerateToken("other-secret", 7, "tenant", time.Hour)
	if err != nil {
		t.Fatalf("GenerateToken returned an error: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	newAuthTestRouter("test-secret").ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}