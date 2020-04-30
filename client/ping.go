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

// assuming ipv4
func Ping(url string) error {
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
	fmt.Println(m)

	packetConn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return err
	}
	defer packetConn.Close()

	mb, err := m.Marshal(nil)
	if err != nil {
		return err
	}

	targetIP := net.ParseIP("1.1.1.1")
	if targetIP == nil {
		return errors.New("invalid ip")
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
		return errors.New(fmt.Sprintf("connection error: %s", err.Error()))
	}

	responseMessage, err := icmp.ParseMessage(ipv4.ICMPTypeEcho.Protocol(), rb[:respSize])
	fmt.Println("respond time: ", duration, "responseType: ", responseMessage.Type)
	return nil
}
