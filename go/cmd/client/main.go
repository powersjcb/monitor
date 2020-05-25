package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/jackpal/gateway"
	"github.com/powersjcb/monitor/go/src/client"
	"github.com/powersjcb/monitor/go/src/lib/tracer"
	"log"
	"os"
	"strings"
	"time"
)

const publicTracingKey = "a0f88ec0416dae30766466ab00f0492c"

func main() {
	runOnce := flag.Bool("once", false, "Run pings only one time.")
	flag.Parse()
	tracer.InitTracer(publicTracingKey)
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
	apiKey := os.Getenv("MONITOR_API_KEY")
	if apiKey == "" {
		fmt.Println("warning invalid api key")
	}
	err = client.RunPings(ctx, apiKey, pingConfigs, *runOnce, h)
	if err != nil {
		log.Fatal(err)
	}
}
