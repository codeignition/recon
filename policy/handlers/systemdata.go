// Copyright 2015 CodeIgnition. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package handlers

import (
	"errors"
	"log"
	"os/user"
	"time"

	"github.com/codeignition/recon/metrics/netstat"
	"github.com/codeignition/recon/metrics/ps"
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
					Time:   time.Now(),
					Policy: p,
					Data:   accumulateSystemData(),
				}
			}
		}
	}()
	return out, nil
}

func accumulateSystemData() map[string]interface{} {
	currentUser, err := user.Current()
	if err != nil {
		log.Println(err)
	}
	psdata, err := ps.CollectData()
	if err != nil {
		log.Println(err)
	}
	nsdata, err := netstat.CollectData()
	if err != nil {
		log.Println(err)
	}
	data := map[string]interface{}{
		"recon_time":         time.Now(),
		"current_user":       currentUser.Username, // if more data is required, use currentUser instead of just the Username field
		"process_statistics": psdata,
		"network_statistics": nsdata,
	}
	return data
}
