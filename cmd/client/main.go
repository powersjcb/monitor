package main

import (
	"fmt"
	"github.com/powersjcb/monitor/client"
	"log"
	"time"
)

func main() {
	pingConfigs := []client.PingConfig{
		{URL: "google.com"},
	}
	ticker := time.NewTicker(1 * time.Second)
	fmt.Println("start")
	for _ = range ticker.C {
		fmt.Println("tick")
		err := client.RunPings(pingConfigs)
		if err != nil {
			log.Fatalln(err.Error())
		}
	}
}
