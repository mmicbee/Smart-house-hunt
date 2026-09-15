package auth

import "testing"

func TestHashPassword(t *testing.T) {
	hashed, err := HashPassword("secret123")
	if err != nil {
		t.Fatalf("HashPassword returned an error: %v", err)
	}

	if hashed == "" {
		t.Fatal("HashPassword returned an empty hash")
	}
	if hashed == "secret123" {
		t.Fatal("password was stored without hashing")
	}
	if !CheckPassword(hashed, "secret123") {
		t.Fatal("correct password should validate against its hash")
	}
	if CheckPassword(hashed, "wrongpass") {
		t.Fatal("wrong password should not validate against the hash")
	}

	again, err := HashPassword("secret123")
	if err != nil {
		t.Fatalf("second HashPassword call failed: %v", err)
	}
	if hashed == again {
		t.Fatal("bcrypt should use a unique salt per hash")
	}
}

func TestCheckPasswordInvalidHash(t *testing.T) {
	if CheckPassword("not-a-bcrypt-hash", "secret123") {
		t.Fatal("an invalid hash should never validate")
	}
}