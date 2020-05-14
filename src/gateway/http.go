package gateway

import (
	"encoding/json"
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
	port string
}

func NewHTTPServer(q *db.Queries, port string) HTTPServer {
	return HTTPServer{
		q: q,
		port: port,
	}
}

func (s *HTTPServer) Start() error {
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/", s.Status)
	serverMux.HandleFunc("/metric", s.Metric)
	server := &http.Server{
		Addr: "0.0.0.0:" + s.port,
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

func(s HTTPServer) Status(rw http.ResponseWriter, r *http.Request) {
	_, _ = rw.Write([]byte("ok"))
	rw.WriteHeader(200)
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
		log.Println(err.Error())
		_, _ = rw.Write([]byte(err.Error()))
		rw.WriteHeader(500)
	}
}
