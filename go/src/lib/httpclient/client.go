package httpclient

import (
	"context"
	"go.opentelemetry.io/otel/plugin/httptrace"
	"io"
	"net/http"
	"time"
)

// exposes same interface as http.Client, but provides tracing
type Client struct {
	client *http.Client
}

func (c *Client) Get(ctx context.Context, url string) (resp *http.Response, err error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	ctx, req = httptrace.W3C(ctx, req)
	httptrace.Inject(ctx, req)
	return c.client.Do(req)
}

func (c *Client) PostWithContext(ctx context.Context, url, contentType string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return nil, err
	}
	ctx, req = httptrace.W3C(ctx, req)
	httptrace.Inject(ctx, req)
	req.Header.Set("Content-Type", contentType)
	return c.client.Do(req)
}

func New(timeout time.Duration) Client {
	return Client{
		client:  &http.Client{
			Timeout:       timeout,
		},
	}
}
