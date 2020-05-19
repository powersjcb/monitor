package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/jackpal/gateway"
	"github.com/powersjcb/monitor/go/src/client"
	"log"
	"os"
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
			fmt.Printf("ipv6 unimplemented: %s", gw.String())
		}
		pingConfigs = append(pingConfigs, client.PingConfig{URL: gw.String(), Name: "defaultGateway", Period: 5 * time.Second})
	} else if err != nil {
		fmt.Printf("unable to discover default gateway: %s", err.Error())
	}
	h, _ := os.Hostname()
	ctx := context.Background()
	err = client.RunPings(ctx, pingConfigs, *runOnce, h)
	if err != nil {
		log.Fatal(err)
	}
}
