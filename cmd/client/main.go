package main

import (
	"github.com/powersjcb/monitor/src/client"
	"github.com/jackpal/gateway"
	"log"
	"strings"
	"time"
)

func main() {
	pingConfigs := []client.PingConfig{
		{
			URL:    "google.com",
			Period: 5 * time.Second,
		},
		{
			URL: "amazon.com",
			Period: 5 * time.Second,
		},
		{
			URL: "ec2.us-west-1.amazonaws.com",
			Period: 5 * time.Second,
		},
		{
			URL: "cloudflare.com",
			Period: 5 * time.Second,
		},
	}
	gw, err := gateway.DiscoverGateway()
	if err == nil && gw != nil && gw.String() != "" {
		if strings.Contains(gw.String(), ":") {
			log.Printf("ipv6 unimplemented: %s", gw.String())
		}
		pingConfigs = append(pingConfigs, client.PingConfig{URL: gw.String(), Name: "defaultGateway", Period: 5 * time.Second})
	} else if err != nil {
		log.Printf("unable to discover default gateway: %s", err.Error())
	}

	err = client.RunPings(pingConfigs)
	if err != nil {
		log.Fatal(err)
	}
}
