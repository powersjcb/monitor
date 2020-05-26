package gateway

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/pat"
	"github.com/powersjcb/monitor/go/src/client"
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

func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s %s %v\n", r.Method, r.URL.Path, time.Now())
	l.handler.ServeHTTP(w, r)
}

func NewLogger(handler http.Handler) *Logger {
	return &Logger{handler}
}

func NewTracer(handler http.Handler, tracer trace.Tracer) http.Handler {
	return othttp.NewHandler(handler, "server", othttp.WithTracer(tracer))
}

type HTTPServer struct {
	appContext *ApplicationContext
	jwtConfig  JWTConfig
	oauthConfig OAUTHConfig
	port       string
	apiKey     string
}

type ApplicationContext struct {
	Querier      db.Querier
	Tracer       otel.Tracer
	Logger       log.Logger
}

func NewHTTPServer(appContext *ApplicationContext, jwtConfig JWTConfig, oauth OAUTHConfig, port, apiKey string) HTTPServer {
	if appContext.Querier == nil {
		panic("no db.Querier")
	}
	return HTTPServer{
		appContext: appContext,
		jwtConfig:  jwtConfig,
		oauthConfig: oauth,
		port:       port,
		apiKey: apiKey,
	}
}

func (s *HTTPServer) Start() error {
	p := pat.New()
	// public endpoints
	p.Get(loginPath, s.GoogleLoginHandler)
	p.Get("/auth/google/callback", s.GoogleCallbackHandler)
	p.Get("/status", s.Status)
	p.Get("/pings", s.Ping)

	// endpoints requiring authorization
	p.Post("/metric", s.Authenticated(s.Metric))
	p.Get("/api/profile", s.Authenticated(s.ShowAPIKey))
	p.Get("/", s.Authenticated(s.ShowAPIKey))

	// start server
	server := &http.Server{
		Addr:         "0.0.0.0:" + s.port,
		Handler:      NewTracer(NewLogger(p), s.appContext.Tracer),
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
		fmt.Println(err.Error())
		rw.WriteHeader(500)
		return
	}

	id, err := UserIDFromContext(r.Context())
	if err != nil {
		rw.WriteHeader(500)
		return
	}
	m.AccountID = id
	_, err = s.appContext.Querier.InsertMetric(r.Context(), m)
	if err != nil {
		fmt.Println(err.Error())
		rw.WriteHeader(500)
		return
	}
}

func (s HTTPServer) Ping(rw http.ResponseWriter, r *http.Request) {
	err := client.RunHTTPPings(r.Context(), s.apiKey, client.DefaultPingConfigs, true, r.Host)
	if err != nil {
		fmt.Printf("failed to run pings: %s", err.Error())
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
