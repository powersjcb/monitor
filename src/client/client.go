package client

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/powersjcb/monitor/src/server/db"
	"net/http"
	"time"
)

type PingConfig struct {
	URL string
}

// with some overall timeout, ping all the services
// todo: return some results
func RunPings(configs []PingConfig) error {
	c := PingClient{Timeout: time.Second * 10}

	ticker := time.NewTicker(5 * time.Second)
	for _ = range ticker.C {
		for _, config := range configs {
			p, err := c.Ping(config.URL)
			if err != nil {
				fmt.Println(config.URL, " ", err)
				continue
			}

			body, err := json.Marshal(&db.InsertMetricParams{
				Ts:     sql.NullTime{Time: p.Timestamp, Valid: true},
				Source: "Jacobs-MacBook-Pro", // this computer's hostname
				Name:   "ping",
				Target: p.Target,
				Value:  sql.NullFloat64{Float64: p.Duration.Seconds(), Valid: true},
			})

			fmt.Println("target: ", p.Target, " duration: ", p.Duration)
			_, err = http.Post("https://carbide-datum-276117.wl.r.appspot.com/metric", "application/json", bytes.NewBuffer(body))
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
	}
	return nil
}
