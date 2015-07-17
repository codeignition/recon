// Copyright 2015 Hari haran. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

// metricsPath in the master server
const metricsPath = "/metrics"

func main() {
	log.SetPrefix("recon: ")

	var masterAddr = flag.String("masterAddr", "http://localhost:3000", "address of the recon-master server (along with protocol)")
	flag.Parse()

	c := time.Tick(5 * time.Second)
	for now := range c {
		log.Println("Update sent at", now)
		if err := update(*masterAddr); err != nil {
			log.Println(err)
		}
	}
}

func update(addr string) error {
	var buf bytes.Buffer
	d := map[string]map[string]interface{}{
		"metric": {
			"data": accumulateData(),
		},
	}
	if err := json.NewEncoder(&buf).Encode(&d); err != nil {
		return err
	}
	resp, err := http.Post(addr+metricsPath, "application/json", &buf)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("response status code not %d; response is %v\n", http.StatusOK, resp)
	}
	return nil
}
