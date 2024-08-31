package auth

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
)

// ResponseCapture wraps an http.ResponseWriter to capture response data.
type ResponseCapture struct {
	http.ResponseWriter
	Body *bytes.Buffer
}

// Write captures the response body data.
func (rc *ResponseCapture) Write(b []byte) (int, error) {
	rc.Body.Write(b)                  // Write to buffer
	return rc.ResponseWriter.Write(b) // Also write to the original ResponseWriter
}

// HashSigner returns a middleware that signs the response body with an HMAC hash.
func HashSigner(key string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if key == "" {
				next.ServeHTTP(w, r)
				return
			}

			// Capture the response
			capture := &ResponseCapture{
				ResponseWriter: w,
				Body:           &bytes.Buffer{},
			}

			// Serve the next handler with the captured response writer
			next.ServeHTTP(capture, r)

			// Compute HMAC of the response body
			h := hmac.New(sha256.New, []byte(key))
			h.Write(capture.Body.Bytes())
			computedHash := h.Sum(nil)

			// Convert hash to hexadecimal string
			hashString := hex.EncodeToString(computedHash)

			// Set the computed hash in the header of the original ResponseWriter
			w.Header().Set("HashSHA256", hashString)
		})
	}
}
