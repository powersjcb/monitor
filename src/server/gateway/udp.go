package gateway

import (
	"fmt"
	"net"
)

type UDPService struct {}

func (u UDPService) Start() error {
	pc, err := net.ListenPacket("udp", ":10001")
	if err != nil {
		return err
	}
	defer pc.Close()

	for {
		pb := make([]byte, 1500)
		_, _, err := pc.ReadFrom(pb)
		if err != nil {
			fmt.Println("error: ", err.Error())
			continue
		}
	}
}
