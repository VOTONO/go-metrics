package compressor

import (
	"testing"
)

func BenchmarkGzipCompress(b *testing.B) {
	data := []byte("The quick brown fox jumps over the lazy dog. The quick brown fox jumps over the lazy dog.")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = GzipCompress(data)
	}
}
