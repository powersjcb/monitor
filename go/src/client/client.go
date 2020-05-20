package client

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/powersjcb/monitor/go/src/lib/httpclient"
	"github.com/powersjcb/monitor/go/src/server/db"
	"google.golang.org/api/googleapi"
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

func (h UploadHandler) Handle(ctx context.Context, result PingResult, err error) error {
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
	defer googleapi.CloseBody(resp)
	return nil
}

func RunPings(ctx context.Context, configs []PingConfig, runOnce bool, source string) error {
	c := NewService(ctx, configs, time.Second*10, runOnce)
	c.AddHandler(LoggingHandler{})
	c.AddHandler(UploadHandler{
		HTTP:      httpclient.New(5 * time.Second),
		Source:    source,
		Kind:      "icmp",
		UploadURL: "http://127.0.0.1:8080/metric",
	})
	return c.Start()
}

func RunHTTPPings(ctx context.Context, configs []PingConfig, runOnce bool, source string) error {
	c := NewHTTPService(ctx, configs, time.Second*1, runOnce)
	c.AddHandler(LoggingHandler{})
	c.AddHandler(UploadHandler{
		HTTP:      httpclient.New(5 * time.Second),
		Source:    source,
		Kind:      "http",
		UploadURL: "https://carbide-datum-276117.wl.r.appspot.com/metric",
	})
	return c.Start()
}
