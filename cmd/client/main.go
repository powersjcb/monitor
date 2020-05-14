package main

import (
	"flag"
	"github.com/jackpal/gateway"
	"github.com/powersjcb/monitor/src/client"
	"log"
	"strings"
	"time"
)

func main() {
	runOnce := flag.Bool("once", false, "Run pings only one time.")
	flag.Parse()
	pingConfigs := client.DefaultPingConfigs
	gw, err := gateway.DiscoverGateway()
	if err == nil && gw != nil && gw.String() != "" {
		if strings.Contains(gw.String(), ":") {
			log.Printf("ipv6 unimplemented: %s", gw.String())
		}
		pingConfigs = append(pingConfigs, client.PingConfig{URL: gw.String(), Name: "defaultGateway", Period: 5 * time.Second})
	} else if err != nil {
		log.Printf("unable to discover default gateway: %s", err.Error())
	}

	err = client.RunPings(pingConfigs, *runOnce)
	if err != nil {
		log.Fatal(err)
	}
}
