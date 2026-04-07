package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

const testSecret = "supersecretkey"

func TestHashAndCheckPassword(t *testing.T) {
	password := "mysecurepassword"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error hashing password, got %v", err)
	}

	ok, err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("expected no error checking password, got %v", err)
	}
	if !ok {
		t.Fatalf("expected password to match")
	}

	// wrong password case
	ok, err = CheckPasswordHash("wrongpassword", hash)
	if err != nil {
		t.Fatalf("expected no error checking wrong password, got %v", err)
	}
	if ok {
		t.Fatalf("expected password NOT to match")
	}
}

func TestMakeAndValidateJWT(t *testing.T) {
	userID := uuid.New()

	token, err := MakeJWT(userID, testSecret, time.Hour)
	if err != nil {
		t.Fatalf("expected no error creating token, got %v", err)
	}

	parsedID, err := ValidateJWT(token, testSecret)
	if err != nil {
		t.Fatalf("expected no error validating token, got %v", err)
	}

	if parsedID != userID {
		t.Fatalf("expected userID %v, got %v", userID, parsedID)
	}
}

func TestValidateJWT_WrongSecret(t *testing.T) {
	userID := uuid.New()

	token, err := MakeJWT(userID, testSecret, time.Hour)
	if err != nil {
		t.Fatalf("error creating token: %v", err)
	}

	_, err = ValidateJWT(token, "wrongsecret")
	if err == nil {
		t.Fatalf("expected error with wrong secret")
	}
}

func TestValidateJWT_ExpiredToken(t *testing.T) {
	userID := uuid.New()

	// token already expired
	token, err := MakeJWT(userID, testSecret, -time.Hour)
	if err != nil {
		t.Fatalf("error creating token: %v", err)
	}

	_, err = ValidateJWT(token, testSecret)
	if err == nil {
		t.Fatalf("expected error for expired token")
	}
}

func TestGetBearerToken(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer my-token-123")

	token, err := GetBearerToken(headers)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if token != "my-token-123" {
		t.Fatalf("expected token 'my-token-123', got %s", token)
	}
}

func TestGetBearerToken_MissingHeader(t *testing.T) {
	headers := http.Header{}

	_, err := GetBearerToken(headers)
	if err == nil {
		t.Fatalf("expected error for missing header")
	}
}

func TestGetBearerToken_InvalidFormat(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "Invalid token")

	_, err := GetBearerToken(headers)
	if err == nil {
		t.Fatalf("expected error for invalid format")
	}
}

func TestGetBearerToken_EmptyToken(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer ")

	_, err := GetBearerToken(headers)
	if err == nil {
		t.Fatalf("expected error for empty token")
	}
}

