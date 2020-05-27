package client

import "time"

var DefaultPingConfigs = []PingConfig{
	{
		URL:    "google.com",
		Period: 5 * time.Second,
	},
	{
		URL:    "amazon.com",
		Period: 5 * time.Second,
	},
	{
		URL:    "ec2.us-west-1.amazonaws.com",
		Period: 5 * time.Second,
	},
	{
		URL:    "cloudflare.com",
		Period: 5 * time.Second,
	},
	{
		URL:    "wl.r.appspot.com",
		Name:   "app-engine@us-west2",
		Period: 5 * time.Second,
	},
	{
		URL:    "proxy.golang.org",
		Period: 30 * time.Second,
	},
	{
		URL:    "github.com",
		Period: 30 * time.Second,
	}}
