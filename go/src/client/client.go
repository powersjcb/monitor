package client

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/powersjcb/monitor/go/src/lib/httpclient"
	"github.com/powersjcb/monitor/go/src/server/db"
	"time"
)

type PingConfig struct {
	Name   string
	URL    string
	Period time.Duration
}

// uploads the data to server
type UploadHandler struct {
	HTTP      httpclient.Client
	Kind      string
	UploadURL string
	Source    string
}

func (h UploadHandler) Handle(ctx context.Context, result PingResult, e error) error {
	if e != nil {
		fmt.Println("failed to get a result")
	}
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
	resp, err := h.HTTP.PostWithContext(ctx, h.UploadURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("failed to upload results", err)
	}
	defer httpclient.CloseBody(resp)
	return nil
}

func RunPings(ctx context.Context, apiKey string, configs []PingConfig, runOnce bool, source string) error {
	c := NewService(ctx, configs, time.Second*10, runOnce)
	c.AddHandler(LoggingHandler{})
	c.AddHandler(UploadHandler{
		HTTP:      httpclient.New(10 * time.Second, apiKey),
		Source:    source,
		Kind:      "icmp",
		UploadURL: "https://monitor.jacobpowers.me/api/metric",
	})
	return c.Start()
}

func RunHTTPPings(ctx context.Context, apiKey string, configs []PingConfig, runOnce bool, source string) error {
	c := NewHTTPService(ctx, configs, time.Second*1, runOnce)
	c.AddHandler(LoggingHandler{})
	c.AddHandler(UploadHandler{
		HTTP:      httpclient.New(10 * time.Second, apiKey),
		Source:    source,
		Kind:      "http",
		UploadURL: "https://monitor.jacobpowers.me/api/metric",
	})
	return c.Start()
}
