package gateway

import (
	"encoding/json"
	"fmt"
	"github.com/powersjcb/monitor/src/server/db"
	"log"
	"net/http"
	"time"
)

type Logger struct {
	handler http.Handler
}

func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	l.handler.ServeHTTP(w, r)
	log.Printf("%s %s %v", r.Method, r.URL.Path, time.Now())
}

func NewLogger(handler http.Handler) *Logger {
	return &Logger{handler}
}

type HTTPServer struct {
	q *db.Queries
}

func NewHTTPServer(q *db.Queries) HTTPServer {
	return HTTPServer{q: q}
}

func (s *HTTPServer) Start() error {
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/metric", s.Metric)
	server := &http.Server{
		Addr: "127.0.0.1:8080",
		Handler: NewLogger(serverMux),
		ReadTimeout: 50 * time.Millisecond,
		WriteTimeout: 50 * time.Millisecond,
	}

	err := server.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (s HTTPServer) Metric(rw http.ResponseWriter, r *http.Request) {
	var m db.InsertMetricParams
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		_, _ = rw.Write([]byte(err.Error()))
		rw.WriteHeader(500)
	}

	_, err = s.q.InsertMetric(r.Context(), m)
	if err != nil {
		fmt.Println(err.Error())
		_, _ = rw.Write([]byte(err.Error()))
		rw.WriteHeader(500)
	}
}
