package auth

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/VOTONO/go-metrics/internal/constants"
)

const (
	testKey  = "test_secret_key"
	testBody = "Hello, world!"
)

func TestHashSigner(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	signedHandler := HashSigner(testKey)(handler)

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	signedHandler.ServeHTTP(rec, req)

	// Compute expected HMAC for comparison
	h := hmac.New(sha256.New, []byte(testKey))
	h.Write([]byte(testBody))
	expectedHash := hex.EncodeToString(h.Sum(nil))

	// Verify the response hash header
	if rec.Header().Get(constants.HashSHA256) != expectedHash {
		t.Errorf("HashSigner() = %v, want %v", rec.Header().Get(constants.HashSHA256), expectedHash)
	}

	// Verify response body content
	if rec.Body.String() != testBody {
		t.Errorf("Expected body = %v, got %v", testBody, rec.Body.String())
	}
}

func TestHashChecker(t *testing.T) {
	// Prepare a request body
	body := "Request body content"
	bodyBytes := []byte(body)

	// Compute HMAC for the request body
	h := hmac.New(sha256.New, []byte(testKey))
	h.Write(bodyBytes)
	hash := hex.EncodeToString(h.Sum(nil))

	// Set up a handler to capture if the request passes the hash check
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(testBody))
	})

	// Wrap the handler with HashChecker middleware
	checkedHandler := HashChecker(testKey)(handler)

	// Create a request with the correct hash in the header
	req := httptest.NewRequest("POST", "/", bytes.NewReader(bodyBytes))
	req.Header.Set(constants.HashSHA256, hash)
	rec := httptest.NewRecorder()

	// Serve the request
	checkedHandler.ServeHTTP(rec, req)

	// Verify the response
	if rec.Code != http.StatusOK {
		t.Errorf("HashChecker() status = %v, want %v", rec.Code, http.StatusOK)
	}
	if rec.Body.String() != testBody {
		t.Errorf("HashChecker() body = %v, want %v", rec.Body.String(), "Hash verified")
	}
}

func TestHashChecker_InvalidHash(t *testing.T) {
	body := "Some different content"
	req := httptest.NewRequest("POST", "/", strings.NewReader(body))
	req.Header.Set(constants.HashSHA256, "invalid_hash")

	rec := httptest.NewRecorder()
	checkedHandler := HashChecker(testKey)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(testBody))
	}))

	// Serve the request
	checkedHandler.ServeHTTP(rec, req)

	// Check for status 400 and error message
	if rec.Code != http.StatusBadRequest {
		t.Errorf("HashChecker() status = %v, want %v", rec.Code, http.StatusBadRequest)
	}
	if !strings.Contains(rec.Body.String(), "Invalid hash") {
		t.Errorf("Expected error message in response body, got %v", rec.Body.String())
	}
}
