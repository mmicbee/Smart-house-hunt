package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"smart-house-hunt/internal/auth"
)

func newRoleTestRouter(secret string, roles ...string) *gin.Engine {
	r := gin.New()
	protected := r.Group("/admin", RequireAuth(secret), RequireRole(roles...))
	protected.GET("/users", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	return r
}

func roleRequest(t *testing.T, secret, role string) *httptest.ResponseRecorder {
	token, err := auth.GenerateToken(secret, 1, role, time.Hour)
	if err != nil {
		t.Fatalf("GenerateToken returned an error: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/admin/users", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	newRoleTestRouter(secret, "admin").ServeHTTP(w, req)
	return w
}

func TestRequireRoleAllowed(t *testing.T) {
	w := roleRequest(t, "test-secret", "admin")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for admin, got %d (%s)", w.Code, w.Body.String())
	}
}

func TestRequireRoleDenied(t *testing.T) {
	for _, role := range []string{"tenant", "landlord", "agent"} {
		w := roleRequest(t, "test-secret", role)
		if w.Code != http.StatusForbidden {
			t.Errorf("expected 403 for %s, got %d", role, w.Code)
		}
	}
}

func TestRequireRoleMultipleAllowed(t *testing.T) {
	token, err := auth.GenerateToken("test-secret", 1, "agent", time.Hour)
	if err != nil {
		t.Fatalf("GenerateToken returned an error: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/admin/users", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	newRoleTestRouter("test-secret", "admin", "agent").ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for agent with multi-role gate, got %d", w.Code)
	}
}