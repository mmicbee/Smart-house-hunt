package auth

import (
	"strings"
	"testing"
	"time"
)

const testSecret = "test-secret"

func TestGenerateAndParseToken(t *testing.T) {
	token, err := GenerateToken(testSecret, 42, "landlord", time.Hour)
	if err != nil {
		t.Fatalf("GenerateToken returned an error: %v", err)
	}

	claims, err := ParseToken(testSecret, token)
	if err != nil {
		t.Fatalf("ParseToken returned an error: %v", err)
	}

	if claims.UserID != 42 {
		t.Errorf("expected user_id 42, got %d", claims.UserID)
	}
	if claims.Role != "landlord" {
		t.Errorf("expected role landlord, got %s", claims.Role)
	}
	if claims.Issuer != "smart-house-hunt" {
		t.Errorf("expected issuer smart-house-hunt, got %s", claims.Issuer)
	}
}

func TestTokenStructure(t *testing.T) {
	token, err := GenerateToken(testSecret, 1, "tenant", time.Hour)
	if err != nil {
		t.Fatalf("GenerateToken returned an error: %v", err)
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Fatalf("expected a 3-part JWT, got %d parts", len(parts))
	}
}

func TestParseTokenWrongSecret(t *testing.T) {
	token, err := GenerateToken(testSecret, 1, "tenant", time.Hour)
	if err != nil {
		t.Fatalf("GenerateToken returned an error: %v", err)
	}

	if _, err := ParseToken("a-different-secret", token); err == nil {
		t.Fatal("token signed with a different secret should be rejected")
	}
}

func TestParseTokenTampered(t *testing.T) {
	token, err := GenerateToken(testSecret, 1, "tenant", time.Hour)
	if err != nil {
		t.Fatalf("GenerateToken returned an error: %v", err)
	}

	parts := strings.Split(token, ".")
	parts[1] = "eyJyb2xlIjoiYWRtaW4ifQ"
	tampered := strings.Join(parts, ".")

	if _, err := ParseToken(testSecret, tampered); err == nil {
		t.Fatal("a tampered token should be rejected")
	}
}

func TestParseTokenExpired(t *testing.T) {
	token, err := GenerateToken(testSecret, 1, "tenant", -time.Minute)
	if err != nil {
		t.Fatalf("GenerateToken returned an error: %v", err)
	}

	if _, err := ParseToken(testSecret, token); err == nil {
		t.Fatal("an expired token should be rejected")
	}
}

func TestParseTokenGarbage(t *testing.T) {
	if _, err := ParseToken(testSecret, "not-a-jwt"); err == nil {
		t.Fatal("garbage input should be rejected")
	}
}