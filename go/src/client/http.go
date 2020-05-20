package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/powersjcb/monitor/go/src/lib/dns"
	"go.opentelemetry.io/otel/plugin/httptrace"
	"golang.org/x/net/icmp"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type HTTPService struct {
	ctx            context.Context
	mux            sync.Mutex
	conn           *icmp.PacketConn
	dnsEntries     map[string]net.IP
	RunOnce        bool
	ResultHandlers []ResultHandler
	resultsCount   uint64
	Inflight       LRU
	Targets        []PingConfig
	Timeout        time.Duration
}

func NewHTTPService(ctx context.Context, targets []PingConfig, timeout time.Duration, runOnce bool) HTTPService {
	if timeout < 1*time.Millisecond {
		panic("timeout too small: " + string(timeout))
	}
	return HTTPService{
		ctx:            ctx,
		mux:            sync.Mutex{},
		dnsEntries:     make(map[string]net.IP),
		RunOnce:        runOnce,
		ResultHandlers: nil,
		Targets:        targets,
		Timeout:        timeout,
	}
}

func (s *HTTPService) AddHandler(handler ResultHandler) {
	s.mux.Lock()
	s.ResultHandlers = append(s.ResultHandlers, handler)
	s.mux.Unlock()
}

func (s *HTTPService) Start() error {
	for true {
		for _, target := range s.Targets {
			if s.RunOnce {
				res, err := s.send(target)
				err = s.evalHandlers(res, err)
				if err != nil {
					fmt.Printf(err.Error())
				}
			} else {
				// add ticker
			}
		}
		if s.RunOnce {
			return nil
		}
	}
	return nil
}

func (s *HTTPService) evalHandlers(r PingResult, err error) error {
	for _, h := range s.ResultHandlers {
		e := h.Handle(s.ctx, r, err)
		if e != nil {
			return errors.New("handler failed: " + e.Error() + err.Error())
		}
	}
	return nil
}

func (s *HTTPService) send(target PingConfig) (PingResult, error) {
	res := PingResult{
		Target:    target.URL,
		Timestamp: time.Now(),
	}
	ip, err := s.dnsLookup(target.URL)
	if err != nil {
		return res, err
	}

	u, err := ensureHTTP(target.URL)
	if err != nil {
		return res, errors.New("invalid url: " + err.Error())
	}

	t := time.Now()
	_, err = get(s.ctx, u, ip, s.Timeout)
	if err != nil {
		return res, err
	}
	res.Duration = time.Now().Sub(t)
	return res, nil
}

func (s *HTTPService) dnsLookup(host string) (net.IP, error) {
	s.mux.Lock()
	e, exists := s.dnsEntries[host]
	s.mux.Unlock()
	if exists {
		return e, nil
	}

	r := dns.NewResolver(s.Timeout)
	ctx := context.Background()
	ips, err := r.LookupIPAddr(ctx, host)
	if err != nil {
		return nil, err
	}

	if len(ips) == 0 {
		return nil, errors.New(fmt.Sprintf("empty list of ips for : %s", host))
	}
	s.mux.Lock()
	s.dnsEntries[host] = ips[0].IP
	s.mux.Unlock()
	return ips[0].IP, nil
}

// load balancers love telling you that you're wrong
// we can abuse fast responses from 301/302s for https redirect
func ensureHTTP(urlString string) (string, error) {
	s := "http://" + urlString
	_, err := url.Parse(s)
	if err != nil {
		return s, err
	}
	return s, nil
}

// we want to measure latency without dns lookup
func get(ctx context.Context, urlString string, cachedIP net.IP, timeout time.Duration) (resp *http.Response, err error) {
	dialer := &net.Dialer{
		Timeout: timeout,
	}

	formattedIP := cachedIP.String()
	if cachedIP.To4() == nil {
		// ivp6 requires [::]:80 formatting
		formattedIP = "[" + formattedIP + "]"
	}

	// Create a transport like http.DefaultTransport, but with a specified localAddr
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: func(ctx context.Context, network, _ string) (net.Conn, error) {
			return dialer.DialContext(ctx, network, formattedIP+":80")
		},
		MaxIdleConns:          100,
		IdleConnTimeout:       timeout,
		TLSHandshakeTimeout:   timeout,
		ExpectContinueTimeout: timeout,
	}
	c := http.Client{
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequestWithContext(ctx, "HEAD", urlString, nil)
	if err != nil {
		return nil, err
	}
	ctx, req = httptrace.W3C(ctx, req)
	httptrace.Inject(ctx, req)

	res, err := c.Do(req)
	defer res.Body.Close()
	return res, err
}
