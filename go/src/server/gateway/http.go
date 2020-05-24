package gateway

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/pat"
	"github.com/powersjcb/monitor/go/src/client"
	"github.com/powersjcb/monitor/go/src/server/db"
	"github.com/powersjcb/monitor/go/src/server/usecases"
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

func NewTracer(handler http.Handler, tracer trace.Tracer) http.Handler {
	return othttp.NewHandler(handler, "server", othttp.WithTracer(tracer))
}

type HTTPServer struct {
	appContext *ApplicationContext
	jwtConfig JWTConfig
	port       string
}

type ApplicationContext struct {
	Querier db.Querier
	Tracer otel.Tracer
	Logger log.Logger
}

func NewHTTPServer(appContext *ApplicationContext, jwtConfig JWTConfig, port string) HTTPServer {
	if appContext.Querier == nil {
		panic("no db.Querier")
	}
	return HTTPServer{
		appContext: appContext,
		port:       port,
		jwtConfig: jwtConfig,
	}
}

func (s *HTTPServer) Start() error {
	p := pat.New()
	p.Post("/metric", s.Metric)
	p.Get("/pings", s.Ping)
	p.Get("/status", s.Status)
	p.Get("/auth/google/login", s.GoogleLoginHandler)
	p.Get("/auth/google/callback", s.GoogleCallbackHandler)

	p.Get("/", s.ShowAPIKey)
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

const googleEmailHeader = "X-Goog-Authenticated-User-Email"
const googleIDHeader = "X-Goog-Authenticated-User-ID"

func (s HTTPServer) ShowAPIKey(rw http.ResponseWriter, r *http.Request) {
	email := r.Header.Get(googleEmailHeader)
	id := r.Header.Get(googleIDHeader)

	account, err := usecases.GetOrCreateAccount(r.Context(), s.appContext.Querier, "google", id)
	if err != nil {
		_, err = rw.Write([]byte(err.Error()))
		if err != nil {
			fmt.Printf("failed to write error response: %s", err.Error())
		}
		rw.WriteHeader(500)
		return
	}

	_, err = rw.Write([]byte(fmt.Sprintf("email: %s, apiKey: %s", email, account.ApiKey)))
	if err != nil {
		_, err = rw.Write([]byte(err.Error()))
		if err != nil {
			fmt.Printf("failed to write error response: %s", err.Error())
		}
		rw.WriteHeader(500)
		return
	}
}

//func (s HTTPServer) staticHandler(rw http.ResponseWriter, r *http.Request) {
//	publicDir := "./public"
//	fi, err := os.Stat(publicDir)
//	if err != nil || !fi.IsDir() {
//		return
//	}
//	staticsMap := map[string]string{
//		"": "./public/index.html",
//		"/": "./public/index.html",
//		"/static/index.js": "public/index.js",
//	}
//	static, exists := staticsMap[r.URL.Path]
//	if !exists {
//		rw.WriteHeader(404)
//		return
//	}
//	http.ServeFile(rw, r, static)
//}

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

	_, err = s.appContext.Querier.InsertMetric(r.Context(), m)
	if err != nil {
		fmt.Println(err.Error())
		rw.WriteHeader(500)
		return
	}
}

func (s HTTPServer) Ping(rw http.ResponseWriter, r *http.Request) {
	err := client.RunHTTPPings(r.Context(), client.DefaultPingConfigs, true, r.Host)
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
