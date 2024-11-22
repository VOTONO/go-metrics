package compressor

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkGzipCompress(b *testing.B) {
	data := []byte("The quick brown fox jumps over the lazy dog. The quick brown fox jumps over the lazy dog.")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = GzipCompress(data)
	}
}

// Test Compressor middleware when Accept-Encoding contains "gzip"
func TestCompressor_WithGzip(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Send a response
		w.Write([]byte("Hello, World"))
	})

	// Wrap the handler with the Compressor middleware
	compressedHandler := Compressor(handler)

	// Create a request with Accept-Encoding: gzip
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	// Record the response
	rr := httptest.NewRecorder()
	compressedHandler.ServeHTTP(rr, req)

	// Check the response headers for Content-Encoding: gzip
	assert.Equal(t, "gzip", rr.Header().Get("Content-Encoding"))

	// Check the response body to ensure it is compressed
	body := rr.Body.Bytes()
	gz, err := gzip.NewReader(bytes.NewReader(body))
	assert.NoError(t, err)

	// Read decompressed body
	decompressedBody, err := io.ReadAll(gz) // Replaced ioutil.ReadAll with io.ReadAll
	assert.NoError(t, err)

	// Assert the decompressed body is the original response
	assert.Equal(t, "Hello, World", string(decompressedBody))
}

// Test Compressor middleware when Accept-Encoding does not contain "gzip"
func TestCompressor_WithoutGzip(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Send a response
		w.Write([]byte("Hello, World"))
	})

	// Wrap the handler with the Compressor middleware
	compressedHandler := Compressor(handler)

	// Create a request without Accept-Encoding header or with a different encoding
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept-Encoding", "identity")

	// Record the response
	rr := httptest.NewRecorder()
	compressedHandler.ServeHTTP(rr, req)

	// Check the response headers for Content-Encoding, should not be gzip
	assert.Empty(t, rr.Header().Get("Content-Encoding"))

	// Check the response body
	assert.Equal(t, "Hello, World", rr.Body.String())
}

// Test Decompressor middleware when Content-Encoding contains "gzip"
func TestDecompressor_WithGzip(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify if the body is decompressed
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		assert.Equal(t, "Hello, World", string(body))
	})

	// Wrap the handler with the Decompressor middleware
	decompressedHandler := Decompressor(handler)

	// Create a compressed request body with gzip
	var compressedBody bytes.Buffer
	gzipWriter := gzip.NewWriter(&compressedBody)
	_, err := gzipWriter.Write([]byte("Hello, World"))
	assert.NoError(t, err)
	gzipWriter.Close()

	// Create a request with Content-Encoding: gzip
	req := httptest.NewRequest("POST", "/", &compressedBody)
	req.Header.Set("Content-Encoding", "gzip")

	// Record the response
	rr := httptest.NewRecorder()
	decompressedHandler.ServeHTTP(rr, req)
}
