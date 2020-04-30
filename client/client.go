package client

import "time"

type PingConfig struct {
	URL string
}

// with some overall timeout, ping all the services
// todo: return some results
func RunPings(configs []PingConfig) error {
	c := PingClient{}

	ticker := time.NewTicker(1 * time.Second)
	for _ = range ticker.C {
		for _, config := range configs {
			err := c.Ping(config.URL)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
