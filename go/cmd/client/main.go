package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/jackpal/gateway"
	"github.com/powersjcb/monitor/go/src/client"
	"github.com/powersjcb/monitor/go/src/lib/tracer"
	"go.opentelemetry.io/otel/api/global"
	"log"
	"os"
	"strings"
	"time"
)

const publicTracingKey = "a0f88ec0416dae30766466ab00f0492c"

func main() {
	ctx := context.Background()
	runOnce := flag.Bool("once", false, "Run pings only one time.")
	flag.Parse()

	tracer.InitTracer(publicTracingKey)
	t := global.Tracer("monitor.client")
	_, span := t.Start(ctx, "client.run")
	defer span.End()

	pingConfigs := client.DefaultPingConfigs
	gw, err := gateway.DiscoverGateway()
	if err == nil && !gw.IsUnspecified() && gw.String() != "" {
		if strings.Contains(gw.String(), ":") {
			fmt.Println("ipv6 unimplemented: ", gw.String())
		}
		pingConfigs = append(pingConfigs, client.PingConfig{URL: gw.String(), Name: "defaultGateway", Period: 5 * time.Second})
	} else if err != nil {
		fmt.Println("unable to discover default gateway: ", err.Error())
	}
	h, _ := os.Hostname()
	apiKey := os.Getenv("MONITOR_API_KEY")
	if apiKey == "" {
		fmt.Println("warning invalid api key")
	}
	err = client.RunPings(ctx, apiKey, pingConfigs, *runOnce, h)
	if err != nil {
		log.Fatal(err)
	}
}
