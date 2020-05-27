package httpserver

import (
	"net/http"
	"sync"
)

type StatusRecorder struct {
	http.ResponseWriter
	Status int
	mux sync.Mutex
}

func (s *StatusRecorder) WriteHeader(statusCode int) {
	s.Status = statusCode
	s.ResponseWriter.WriteHeader(statusCode)
}

func NewStatusRecorder(w http.ResponseWriter) *StatusRecorder {
	return &StatusRecorder{
		ResponseWriter: w,
		Status: 0,
	}
}