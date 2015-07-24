// Copyright 2015 CodeIgnition. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"os/user"
	"time"

	"github.com/codeignition/recon/netstat"
	"github.com/codeignition/recon/ps"
)

func copyMap(from, to map[string]interface{}) {
	for k, v := range from {
		to[k] = v
	}
}

// accumulateData accumulates data from all other packages.
func accumulateData() map[string]interface{} {
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
