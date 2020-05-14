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
	UploadURL string
	Timeout time.Duration
}

func (h UploadHandler) Handle(result PingResult, err error) error {
	body, err := json.Marshal(&db.InsertMetricParams{
		Ts:     sql.NullTime{Time: result.Timestamp, Valid: true},
		Source: "Jacobs-MacBook-Pro", // this computer's hostname
		Name:   "ping",
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

func RunPings(configs []PingConfig) error {
	c := NewService(configs, time.Second * 10, true)

	c.AddHandler(LoggingHandler{})
	c.AddHandler(UploadHandler{
		UploadURL: "https://carbide-datum-276117.wl.r.appspot.com/metric",
		Timeout:   time.Second * 5,
	})
	return c.Start()
}
