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
			if key == "" {
				next.ServeHTTP(w, r)
				return
			}

			hash := r.Header.Get("HashSHA256")
			if hash == "" {
				next.ServeHTTP(w, r)
				return
			}

			hashData, err := hex.DecodeString(hash)
			if err != nil {
				http.Error(w, "Invalid hash format", http.StatusBadRequest)
				return
			}

			bodyData, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Failed to read request body", http.StatusInternalServerError)
				return
			}
			r.Body = io.NopCloser(io.MultiReader(bytes.NewReader(bodyData)))

			h := hmac.New(sha256.New, []byte(key))
			h.Write(bodyData)
			computedHash := h.Sum(nil)

			if hmac.Equal(computedHash, hashData) {
				next.ServeHTTP(w, r)
			} else {
				http.Error(w, "Invalid hash", http.StatusBadRequest)
			}
		})
	}
}
