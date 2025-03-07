package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/toxic-development/sysmanix/utils"
)

// Recovery middleware handles panic recovery
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %v\n%s", err, debug.Stack())
				utils.WriteErrorResponse(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// Logger middleware logs incoming requests
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &statusWriter{ResponseWriter: w}

		next.ServeHTTP(sw, r)

		// Ensure we have a status code (default to 200 OK if not set)
		if sw.status == 0 {
			sw.status = http.StatusOK
		}

		log.Printf(
			"%s %s %s %d %s",
			r.RemoteAddr,
			r.Method,
			r.URL,
			sw.status,
			time.Since(start),
		)
	})
}

type statusWriter struct {
	http.ResponseWriter
	status  int
	written bool
}

func (w *statusWriter) WriteHeader(status int) {
	// Only call the underlying WriteHeader once
	if !w.written {
		w.status = status
		w.ResponseWriter.WriteHeader(status)
		w.written = true
	}
}

func (w *statusWriter) Write(b []byte) (int, error) {
	// If WriteHeader wasn't explicitly called before, call it with 200 OK
	if !w.written {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(b)
}
