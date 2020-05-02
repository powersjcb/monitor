package main

import (
	"github.com/powersjcb/monitor/src/client"
	"log"
)

func main() {
	pingConfigs := []client.PingConfig{
		{URL: "google.com"},
		{URL: "amazon.com"},
		{URL: "ec2.us-west-1.amazonaws.com"},
		{URL: "cloudflare.com"},
	}

	err := client.RunPings(pingConfigs)
	if err != nil {
		log.Fatal(err)
	}
}
