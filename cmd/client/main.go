package main

import (
	"github.com/powersjcb/monitor/src/client"
	"github.com/jackpal/gateway"
	"log"
	"strings"
)

func main() {
	pingConfigs := []client.PingConfig{
		{URL: "google.com"},
		{URL: "amazon.com"},
		{URL: "ec2.us-west-1.amazonaws.com"},
		{URL: "cloudflare.com"},
	}

	gw, err := gateway.DiscoverGateway()
	if err == nil && gw != nil && gw.String() != "" {
		if strings.Contains(gw.String(), ":") {
			log.Printf("ipv6 unimplemented: %s", gw.String())
		}
		pingConfigs = append(pingConfigs, client.PingConfig{URL: gw.String(), Name: "defaultGateway"})
	} else if err != nil {
		log.Printf("unable to discover default gateway: %s", err.Error())
	}

	err = client.RunPings(pingConfigs)
	if err != nil {
		log.Fatal(err)
	}
}
