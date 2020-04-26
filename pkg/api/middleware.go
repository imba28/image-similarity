package api

import (
	"log"
	"net/http"
)

type responseLogger struct {
	statusCode int
	w          http.ResponseWriter
}

func (r *responseLogger) ResponseCode() int {
	return r.statusCode
}

func (r *responseLogger) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.w.WriteHeader(statusCode)
}

func (r *responseLogger) Header() http.Header {
	return r.w.Header()
}

func (r *responseLogger) Write(b []byte) (int, error) {
	return r.w.Write(b)
}

func newResponseLogger(w http.ResponseWriter) responseLogger {
	return responseLogger{
		w: w,
	}
}

func LoggingMiddleware(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := newResponseLogger(w)
		next.ServeHTTP(&logger, r)
		log.Printf("%s - \"%s %s\" %d", r.RemoteAddr, r.Method, r.URL, logger.statusCode)
	})
}

var _ http.ResponseWriter = (*responseLogger)(nil)
