package main

import (
	"github.com/powersjcb/monitor/client"
	"log"
)

func main() {
	pingConfigs := []client.PingConfig{
		{URL: "google.com"},
		{URL: "amazon.com"},
		{URL: "cloudflare.com"},
	}

	err := client.RunPings(pingConfigs)
	if err != nil {
		log.Fatal(err)
	}
}
