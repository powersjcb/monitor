package gateway

import (
	"encoding/json"
	"fmt"
	"github.com/powersjcb/monitor/src/client"
	"github.com/powersjcb/monitor/src/server/db"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/plugin/othttp"
	"log"
	"net/http"
	"time"
)

type Logger struct {
	handler http.Handler
}

func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s %s %v\n", r.Method, r.URL.Path, time.Now())
	l.handler.ServeHTTP(w, r)
}

func NewLogger(handler http.Handler) *Logger {
	return &Logger{handler}
}

func NewTracer(handler http.Handler, tracer trace.Tracer)  http.Handler {
	return othttp.NewHandler(handler, "server", othttp.WithTracer(tracer))
}

type HTTPServer struct {
	appContext *ApplicationContext
	port string
	q *db.Queries
}

type ApplicationContext struct {
	Tracer otel.Tracer
	Logger log.Logger
}

func NewHTTPServer(appContext *ApplicationContext, q *db.Queries, port string) HTTPServer {
	return HTTPServer{
		appContext: appContext,
		port: port,
		q: q,
	}
}

func (s *HTTPServer) Start() error {
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/metric", s.Metric)
	serverMux.HandleFunc("/pings", s.Ping)
	serverMux.HandleFunc("/status", s.Status)
	server := &http.Server{
		Addr:         "0.0.0.0:" + s.port,
		Handler:      NewTracer(NewLogger(serverMux), s.appContext.Tracer),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	err := server.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func(s HTTPServer) Status(rw http.ResponseWriter, r *http.Request) {
	_, _ = rw.Write([]byte("ok"))
}

func (s HTTPServer) Metric(rw http.ResponseWriter, r *http.Request) {
	var m db.InsertMetricParams
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		fmt.Printf(err.Error())
		rw.WriteHeader(500)
		return
	}

	_, err = s.q.InsertMetric(r.Context(), m)
	if err != nil {
		fmt.Printf(err.Error())
		rw.WriteHeader(500)
		return
	}
}

func (s HTTPServer) Ping(rw http.ResponseWriter, r *http.Request) {
	err := client.RunHTTPPings(r.Context(), client.DefaultPingConfigs, true, r.Host)
	if err != nil {
		fmt.Printf("failed to run pings", err.Error())
		rw.WriteHeader(500)
		return
	}
	_, err = rw.Write([]byte("pong"))
	if err != nil {
		fmt.Println(err.Error())
		rw.WriteHeader(500)
		return
	}
}