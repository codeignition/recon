// Copyright 2015 CodeIgnition. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package handlers

import (
	"errors"
	"net"
	"time"

	"github.com/codeignition/recon/policy"

	"golang.org/x/net/context"
)

func TCP(ctx context.Context, p policy.Policy) (<-chan policy.Event, error) {
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

	out := make(chan policy.Event)
	go func() {
		for now := range time.Tick(d) {
			_, err := net.DialTimeout("tcp", port, d)
			if err != nil {
				out <- policy.Event{
					Time:   now,
					Policy: p,
					Data: map[string]interface{}{
						"status": "failure",
						"error":  err,
					},
				}
			} else {
				out <- policy.Event{
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
