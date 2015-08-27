package policy

import (
	"errors"
	"net"
	"time"

	"golang.org/x/net/context"
)

func tcpPolicyHandler(ctx context.Context, p Policy) (<-chan Event, error) {
	// Always use v, ok := p[key] form to avoid panic
	port, ok := p.M["port"]
	if !ok {
		return nil, errors.New(`"port" key missing in tcp policy`)
	}
	freq, ok := p.M["frequency"]
	if !ok {
		return nil, errors.New(`"frequency" key missing in tcp policy`)
	}

	// From the time package docs:
	//
	// ParseDuration parses a duration string.
	// A duration string is a possibly signed sequence of
	// decimal numbers, each with optional fraction and a unit suffix,
	// such as "300ms", "-1.5h" or "2h45m".
	// Valid time units are "ns", "us", "ms", "s", "m", "h".
	d, err := time.ParseDuration(freq)
	if err != nil {
		return nil, err
	}

	// This check is here to ensure time.Ticker(d) doesn't panic
	if d <= 0 {
		return nil, errors.New("frequency must be a positive quantity")
	}

	out := make(chan Event)
	go func() {
		for now := range time.Tick(d) {
			_, err := net.DialTimeout("tcp", port, d)
			if err != nil {
				out <- Event{
					Time:   now,
					Policy: p,
					Data: map[string]interface{}{
						"status": "failure",
						"error":  err,
					},
				}
			} else {
				out <- Event{
					Time:   now,
					Policy: p,
					Data: map[string]interface{}{
						"status": "success",
					},
				}
			}
			// TODO: think about when to close the out channel .
		}
	}()
	return out, nil
}
