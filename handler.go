package http_brotli_handler

import (
	"github.com/google/brotli/go/cbrotli"
	"io"
	"net/http"
	"strings"
)

const MaxCompressionLevel = 11

type responseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w *responseWriter) Write(b []byte) (int, error) {
	h := w.ResponseWriter.Header()
	if h.Get("Content-Type") == "" {
		h.Set("Content-Type", http.DetectContentType(b))
	}

	return w.Writer.Write(b)
}

func CompressHandlerLevel(h http.Handler, level int) http.Handler {
	if level < 0 {
		level = 0
	}

	if level > MaxCompressionLevel {
		level = MaxCompressionLevel
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isSupported(r) {
			writer := cbrotli.NewWriter(w, cbrotli.WriterOptions{Quality: level})
			defer writer.Close()

			w.Header().Set("Content-Encoding", "br")
			w.Header().Add("Vary", "Accept-Encoding")

			w = &responseWriter{
				Writer:         writer,
				ResponseWriter: w,
			}
		}

		h.ServeHTTP(w, r)
	})
}

func isSupported(r *http.Request) bool {
	for _, item := range strings.Split(strings.ToLower(r.Header.Get("Accept-Encoding")), ",") {
		if strings.TrimSpace(item) == "br" {
			return true
		}
	}

	return false
}
