package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/powersjcb/monitor/src/lib/dns"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"
)

type PingID uint16

type PingRequest struct {
	ID     PingID
	Seq    uint16
	Target string
	SentAt time.Time
}

// todo return information about success/failure
type PingResult struct {
	Target string
	Duration time.Duration
	Timestamp time.Time
}

type ResultHandler interface {
	Handle(result PingResult, e error) error
}

type LoggingHandler struct {}

func (LoggingHandler) Handle(r PingResult, err error) error {
	if err != nil {
		log.Printf("failed for target %s: %s", r.Target, err.Error())
	} else {
		log.Printf("name: %s, latency: %s\n", r.Target, r.Duration)
	}
	return nil
}

type PingService struct {
	mux sync.Mutex
	conn *icmp.PacketConn
	dnsEntries map[string]net.IP
	RunOnce bool
	ResultHandlers []ResultHandler
	resultsCount uint64
	Inflight LRU
	Targets []PingConfig
	Timeout time.Duration
}

func NewService(targets []PingConfig, timeout time.Duration, runOnce bool) PingService {
	if timeout < 1 * time.Millisecond {
		panic("timeout too small: " + string(timeout))
	}
	return PingService{
		mux:            sync.Mutex{},
		conn:           nil,
		dnsEntries:     make(map[string]net.IP),
		RunOnce: 		runOnce,
		resultsCount:   0,
		ResultHandlers: nil,
		Inflight:       NewLRU(25),
		Targets:        targets,
		Timeout:        timeout,
	}
}

func (c *PingService) getTimeout() time.Duration {
	c.mux.Lock()
	t := c.Timeout
	c.mux.Unlock()
	return t
}

func (c *PingService) getResultsCount() uint64 {
	c.mux.Lock()
	val := c.resultsCount
	c.mux.Unlock()
	return val
}

func (c *PingService) dnsLookup(host string) (net.IP, error) {
	c.mux.Lock()
	e, exists := c.dnsEntries[host]
	c.mux.Unlock()
	if exists {
		return e, nil
	}

	r := dns.NewResolver(c.getTimeout())
	ctx := context.Background()
	ips, err := r.LookupIPAddr(ctx, host)
	if err != nil {
		return nil, err
	}

	if len(ips) == 0 {
		return nil, errors.New(fmt.Sprintf("empty list of ips for : %s", host))
	}
	c.mux.Lock()
	c.dnsEntries[host] = ips[0].IP
	c.mux.Unlock()
	return ips[0].IP, nil
}

func (c *PingService) sendRequest(host string) error {
	msgID := rand.Intn(1 << 15 + 1)
	m := icmp.Message{
		Type:     ipv4.ICMPTypeEcho,
		Code:     0,
		Checksum: 0, // checksum populated by Marshal func
		Body:     &icmp.Echo{
			ID:   msgID, // 16bit number
			Seq:  0,
			Data: nil,
		},
	}
	mb, err := m.Marshal(nil)
	if err != nil {
		return err
	}
	targetIP, err := c.dnsLookup(host)
	if err != nil {
		return errors.New("dns error: " + err.Error())
	}

	c.mux.Lock() // todo: add connection pool
	s := time.Now()
	err = c.conn.SetWriteDeadline(s.Add(c.Timeout))
	c.mux.Unlock()
	if err != nil {
		return err
	}
	_, err = c.conn.WriteTo(
		mb,
		&net.IPAddr{IP: targetIP},
	)
	if err != nil {
		return errors.New("failed to write request: " + err.Error())
	}
	c.mux.Lock()
	c.Inflight.Add(PingRequest{
		ID: PingID(msgID),
		Seq: 0,
		Target: host,
		SentAt: s,
	})
	c.mux.Unlock()
	return nil
}

func (c *PingService) init() error {
	packetConn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return err
	}
	c.conn = packetConn
	if c.RunOnce {
		err = packetConn.SetDeadline(time.Now().Add(c.Timeout))
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *PingService) nextICMP() (*icmp.Echo, error) {
	rb := make([]byte, 1500)
	respSize, _, err := c.conn.ReadFrom(rb)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("connection error: %s", err.Error()))
	}
	resp, err := icmp.ParseMessage(ipv4.ICMPTypeEcho.Protocol(), rb[:respSize])
	if err != nil {
		return nil, err
	}
	echo, ok := resp.Body.(*icmp.Echo)
	if !ok {
		return nil, errors.New("invalid ping response")
	}
	if echo.ID == 0 {
		return echo, errors.New("invalid echo response, ID unset")
	}
	return echo, nil
}

func (c *PingService) listen(wg *sync.WaitGroup) error {
	wg.Done()
	for true {
		icmpMessage, err := c.nextICMP()
		if err != nil {
			err = c.evalHandlers(PingResult{}, err)
			if err != nil {
				return err
			}
		} else {
			p := time.Now()
			c.mux.Lock()
			req := c.Inflight.Remove(PingID(icmpMessage.ID))
			c.mux.Unlock()
			if req == nil {
				err = errors.New("response to unknown request: " + string(icmpMessage.ID))
				err = c.evalHandlers(PingResult{}, err)
				if err != nil {
					return err
				}
				continue
			}
			res := PingResult{
				Target:    req.Target,
				Duration:  p.Sub(req.SentAt),
				Timestamp: time.Now(),
			}
			err = c.evalHandlers(res, err)
			if err != nil {
				return err
			}

			c.mux.Lock()
			c.resultsCount++
			c.mux.Unlock()

			// maybe shutdown
			if c.RunOnce && c.getResultsCount() > 0 {
				c.mux.Lock()
				remaining := c.Inflight.Len()
				c.mux.Unlock()
				if remaining == 0 {
					return nil
				}
			}
		}
	}
	return nil
}

func (c *PingService) evalHandlers(r PingResult, err error) error {
	for _, h := range c.ResultHandlers {
		e := h.Handle(r, err)
		if e != nil {
			return errors.New("handler failed: " + e.Error() + err.Error())
		}
	}
	return nil
}

func (c *PingService) AddHandler(handler ResultHandler) {
	c.mux.Lock()
	c.ResultHandlers = append(c.ResultHandlers, handler)
	c.mux.Unlock()
}

// start server running with pings
// needs to control ticker/timer logic
func (c *PingService) Start() error {
	err := c.init()
	if err != nil {
		return err
	}
	startupWg := &sync.WaitGroup{}
	startupWg.Add(1)

	for _, target := range c.Targets {
		// spawn a ticker for each config
		p := target.Period
		url := target.URL
		go func () {
			startupWg.Wait()
			if c.RunOnce {
				err := c.sendRequest(url)
				if err != nil {
					log.Println(err.Error())
				}
			} else {
				ticker := time.NewTicker(p)
				for _ = range ticker.C {
					err := c.sendRequest(url)
					if err != nil {
						log.Println(err.Error())
					}
				}
			}
		}()
	}

	err = c.listen(startupWg)
	if err != nil {
		return err
	}
	defer func() {
		c.mux.Lock()
		c.conn.Close()
		c.mux.Unlock()
	}()
	return nil
}
