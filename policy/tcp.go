package policy

import (
	"errors"
	"fmt"
	"net"
	"time"
)

func tcpPolicyHandler(p Policy) error {
	// Always use v, ok := p[key] form to avoid panic
	port, ok := p.M["port"]
	if !ok {
		return errors.New(`"port" key missing in tcp policy`)
	}
	freq, ok := p.M["frequency"]
	if !ok {
		return errors.New(`"frequency" key missing in tcp policy`)
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
		return err
	}

	// This check is here to ensure time.Ticker(d) doesn't panic
	if d <= 0 {
		return errors.New("frequency must be a positive quantity")
	}
	go func() {
		for now := range time.Tick(d) {
			_, err := net.DialTimeout("tcp", port, d)
			if err != nil {
				fmt.Println(now, p.Name, "failure")
				// TODO: sendErrorToMarksman(Agent, Policy, err)
			} else {
				fmt.Println(now, p.Name, "success")
				// TODO: sendSuccessToMarksman(Agent, Policy, err)
			}

		}
	}()
	return nil
}
