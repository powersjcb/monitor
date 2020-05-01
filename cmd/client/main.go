package main

import (
	"github.com/powersjcb/monitor/src/client"
	"log"
)

func main() {
	pingConfigs := []client.PingConfig{
		//{URL: "192.168.7.1"},
		{URL: "google.com"},
		{URL: "amazon.com"},
		{URL: "cloudflare.com"},
		{URL: "zoom.us"},
	}

	err := client.RunPings(pingConfigs)
	if err != nil {
		log.Fatal(err)
	}
}
