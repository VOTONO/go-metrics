package auth

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
)

type ResponseCapture struct {
	http.ResponseWriter
	Body *bytes.Buffer
}

func HashSigner(key string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if key == "" {
				next.ServeHTTP(w, r)
				return
			}
			capture := &ResponseCapture{
				ResponseWriter: w,
				Body:           &bytes.Buffer{},
			}

			next.ServeHTTP(capture, r)

			h := hmac.New(sha256.New, []byte(key))
			h.Write(capture.Body.Bytes())
			computedHash := h.Sum(nil)

			hashString := hex.EncodeToString(computedHash)
			w.Header().Set("HashSHA256", hashString)
		})

	}
}
