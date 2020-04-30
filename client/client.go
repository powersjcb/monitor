package client

type PingConfig struct {
	URL string
}

// with some overall timeout, ping all the services
// todo: return some results
func RunPings(configs []PingConfig) error {
	for _, config := range configs {
		err := Ping(config.URL)
		if err != nil {
			return err
		}
	}
	return nil
}
