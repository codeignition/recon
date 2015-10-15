// Copyright 2015 CodeIgnition. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package handlers

import (
	"errors"
	"log"
	"time"

	"github.com/codeignition/recon/metrics/top"
	"github.com/codeignition/recon/policy"
	"golang.org/x/net/context"
)

func SystemData(ctx context.Context, p policy.Policy) (<-chan policy.Event, error) {
	interval, ok := p.M["interval"]
	if !ok {
		return nil, errors.New(`"interval" key missing in systemdata policy`)
	}
	d, err := time.ParseDuration(interval)
	if err != nil {
		return nil, err
	}
	// This check is here to ensure time.Ticker(d) doesn't panic
	if d <= 0 {
		return nil, errors.New("interval must be a positive quantity")
	}

	out := make(chan policy.Event)
	go func() {
		t := time.NewTicker(d)
		for {
			select {
			case <-ctx.Done():
				t.Stop()
				close(out)
				return
			case <-t.C:
				out <- policy.Event{
					Time:       time.Now(),
					PolicyName: p.Name,
					AgentUID:   p.AgentUID,
					Data:       accumulateSystemData(),
				}
			}
		}
	}()
	return out, nil
}

func accumulateSystemData() interface{} {
	d, err := top.CollectData()
	if err != nil {
		log.Print(err)
	}
	return d
}
