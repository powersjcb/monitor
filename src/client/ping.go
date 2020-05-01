package client

import (
	"errors"
	"fmt"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"math/rand"
	"net"
	"time"
)

type PingResult struct {
	Target string
	Duration time.Duration
	Timestamp time.Time
}

type PingClient struct {
	dnsEntries map[string] net.IP
}

func (c PingClient) dnsLookup(host string) (net.IP, error) {
	e, exists := c.dnsEntries[host]
	if exists {
		return e, nil
	}
	ips, err := net.LookupIP(host)
	if err != nil {
		return nil, err
	}
	if len(ips) == 0 {
		return nil, errors.New(fmt.Sprintf("empty list of ips for : %s", host))
	}
	return ips[0], nil
}

func (c PingClient) Ping(host string) (PingResult, error) {
	m := icmp.Message{
		Type:     ipv4.ICMPTypeEcho,
		Code:     0,
		Checksum: 0, // checksum populated by Marshal func
		Body:     &icmp.Echo{
			ID:   rand.Intn(1 << 15 + 1), // 16bit number
			Seq:  0,
			Data: nil,
		},
	}
	targetIP, err := c.dnsLookup(host)
	if err != nil {
		return PingResult{}, err
	}

	packetConn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return PingResult{}, err
	}
	defer packetConn.Close()

	mb, err := m.Marshal(nil)
	if err != nil {
		return PingResult{}, err
	}

	s := time.Now()
	_, err = packetConn.WriteTo(
		mb,
		&net.IPAddr{IP: targetIP},
	)

	rb := make([]byte, 1500)
	respSize, _, err := packetConn.ReadFrom(rb)
	duration := time.Now().Sub(s)
	if err != nil {
		return PingResult{}, errors.New(fmt.Sprintf("connection error: %s", err.Error()))
	}
	_, err = icmp.ParseMessage(ipv4.ICMPTypeEcho.Protocol(), rb[:respSize])
	return PingResult{Target: host, Duration: duration, Timestamp: s}, nil
}
