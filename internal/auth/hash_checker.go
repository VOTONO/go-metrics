package auth

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
)

func HashChecker(key string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if key == "" || r.Header.Get("HashSHA256") == "" {
				next.ServeHTTP(w, r)
				return
			}

			hashData, err := hex.DecodeString(r.Header.Get("HashSHA256"))
			if err != nil {
				http.Error(w, "Invalid hash format", http.StatusBadRequest)
				return
			}

			h := hmac.New(sha256.New, []byte(key))
			if _, err := io.Copy(h, io.LimitReader(r.Body, 1<<20)); err != nil {
				http.Error(w, "Failed to read request body", http.StatusInternalServerError)
				return
			}
			computedHash := h.Sum(nil)

			r.Body = io.NopCloser(bytes.NewReader(computedHash))

			if hmac.Equal(computedHash, hashData) {
				next.ServeHTTP(w, r)
			} else {
				http.Error(w, "Invalid hash", http.StatusBadRequest)
			}
		})
	}
}
