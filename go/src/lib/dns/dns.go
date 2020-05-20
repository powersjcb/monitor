package dns

import (
	"context"
	"net"
	"time"
)

func NewResolver(timeout time.Duration) net.Resolver {
	return net.Resolver{
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: timeout,
			}
			return d.DialContext(ctx, "udp", "1.1.1.1:53")
		},
	}
}
