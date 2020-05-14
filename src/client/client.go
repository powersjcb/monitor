package client

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/powersjcb/monitor/src/server/db"
	"log"
	"net/http"
	"time"
)

type PingConfig struct {
	Name string
	URL string
	Period time.Duration
}

// uploads the data to server
type UploadHandler struct {
	Kind 	  string
	UploadURL string
	Timeout   time.Duration
	Source 	  string
}

func (h UploadHandler) Handle(result PingResult, err error) error {
	body, err := json.Marshal(&db.InsertMetricParams{
		Ts:     sql.NullTime{Time: result.Timestamp, Valid: true},
		Source: h.Source, // this computer's hostname
		Name:   h.Kind,
		Target: result.Target,
		Value:  sql.NullFloat64{Float64: result.Duration.Seconds(), Valid: true},
	})
	if err != nil {
		return err
	}
	_, err = http.Post(h.UploadURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Println("failed to upload results", err)
	}
	return nil
}

func RunPings(configs []PingConfig, runOnce bool) error {
	c := NewService(configs, time.Second * 10, runOnce)
	c.AddHandler(LoggingHandler{})
	c.AddHandler(UploadHandler{
		Source: "Jacobs-MacBook-Pro.local",
		Kind: "icmp",
		UploadURL: "https://carbide-datum-276117.wl.r.appspot.com/metric",
		Timeout:   time.Second * 5,
	})
	return c.Start()
}

func RunHTTPPings(configs []PingConfig, runOnce bool, source string) error {
	c := NewHTTPService(configs, time.Second * 5, runOnce)
	c.AddHandler(LoggingHandler{})
	c.AddHandler(UploadHandler{
		Source: source,
		Kind: "http",
		UploadURL: "https://carbide-datum-276117.wl.r.appspot.com/metric",
		Timeout:   time.Second * 5,
	})
	return c.Start()
}
