package middleware

import (
	"net/http"
	"sync"
)

// ref: https://upgear.io/blog/golang-tip-wrapping-http-response-writer-for-middleware/
type StatusRecorder struct {
	http.ResponseWriter
	Status int
	mux    sync.Mutex
}

func (s *StatusRecorder) WriteHeader(statusCode int) {
	s.Status = statusCode
	s.ResponseWriter.WriteHeader(statusCode)
}

func NewStatusRecorder(w http.ResponseWriter) *StatusRecorder {
	return &StatusRecorder{
		ResponseWriter: w,
		Status:         0,
	}
}
