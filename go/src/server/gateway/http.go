package gateway

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/powersjcb/monitor/go/src/client"
	"github.com/powersjcb/monitor/go/src/lib/middleware"
	"github.com/powersjcb/monitor/go/src/server/db"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/plugin/othttp"
	"log"
	"net/http"
	"time"
)

const loginPath = "/auth/google/login"

type Logger struct {
	handler http.Handler
}

func (l *Logger) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	w := middleware.NewStatusRecorder(rw)
	start := time.Now()
	// request
	fmt.Printf("request - method: %s path: %s %v\n", r.Method, r.URL.Path, start)
	l.handler.ServeHTTP(w, r)
	// response
	fmt.Printf("response - method: %s path: %s status: %d duration: %s\n", r.Method, r.URL.Path, w.Status, time.Now().Sub(start))
}

func NewLogger(handler http.Handler) *Logger {
	return &Logger{handler}
}

func NewTracer(handler http.Handler, tracer trace.Tracer) http.Handler {
	return othttp.NewHandler(handler, "server", othttp.WithTracer(tracer))
}

type HTTPServer struct {
	appContext  *ApplicationContext
	jwtConfig   JWTConfig
	oauthConfig OAUTHConfig
	port        string
	apiKey      string
}

type ApplicationContext struct {
	Querier db.Querier
	Tracer  otel.Tracer
	Logger  log.Logger
}

func NewHTTPServer(appContext *ApplicationContext, jwtConfig JWTConfig, oauth OAUTHConfig, port, apiKey string) HTTPServer {
	if appContext.Querier == nil {
		panic("no db.Querier")
	}
	return HTTPServer{
		appContext:  appContext,
		jwtConfig:   jwtConfig,
		oauthConfig: oauth,
		port:        port,
		apiKey:      apiKey,
	}
}

const (
	get  = "GET"
	post = "POST"
)

func (s *HTTPServer) Start() error {
	r := mux.NewRouter()
	r.HandleFunc(loginPath, s.GoogleLoginHandler).Methods(get)
	r.HandleFunc("/auth/google/callback", s.GoogleCallbackHandler).Methods(get)
	r.HandleFunc("/status", s.Status).Methods(get)
	r.HandleFunc("/pings", s.Ping).Methods(get)

	r.HandleFunc("/api/metric", s.Authenticated(s.Metric)).Methods(post)
	r.HandleFunc("/api/metric/stats", s.Authenticated(s.MetricStats)).Methods(post)
	r.HandleFunc("/api/profile", s.Authenticated(s.ShowAPIKey)).Methods(get)
	r.HandleFunc("/metric", s.Authenticated(s.Metric)).Methods(post) // do not remove until after updating all clients
	r.HandleFunc("/", s.Authenticated(s.ShowAPIKey))

	// start server
	server := &http.Server{
		Addr:         "0.0.0.0:" + s.port,
		Handler:      NewTracer(NewLogger(r), s.appContext.Tracer),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	err := server.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

func (s HTTPServer) Status(rw http.ResponseWriter, _ *http.Request) {
	_, _ = rw.Write([]byte("ok"))
}

func (s HTTPServer) Metric(rw http.ResponseWriter, r *http.Request) {
	var m db.InsertMetricParams
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		fmt.Println("json.Decode: ", err.Error())
		rw.WriteHeader(500)
		return
	}

	id, err := AccountIDFromContext(r.Context())
	if err != nil {
		rw.WriteHeader(500)
		return
	}
	m.AccountID = id
	_, err = s.appContext.Querier.InsertMetric(r.Context(), m)
	if err != nil {
		fmt.Println("InsertMetric:", err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s HTTPServer) Ping(rw http.ResponseWriter, r *http.Request) {
	err := client.RunHTTPPings(r.Context(), s.apiKey, client.DefaultPingConfigs, true, r.Host)
	if err != nil {
		fmt.Println("failed to run pings: ", err.Error())
		rw.WriteHeader(500)
		return
	}
	_, err = rw.Write([]byte("pong"))
	if err != nil {
		fmt.Println("rw.Write: ", err.Error())
		rw.WriteHeader(500)
		return
	}
}
