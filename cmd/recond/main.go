// Copyright 2015 CodeIgnition. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/codeignition/recon/cmd/recond/config"
	"github.com/nats-io/nats"
)

const (
	metricsAPIPath = "/api/metrics" // metrics path in the master server
	agentsAPIPath  = "/api/agents"  // agents path in the master server
)

// natsEncConn is the opened with the URL obtained from marksman.
// It is populated if the agent registers successfully.
var natsEncConn *nats.EncodedConn

func main() {
	log.SetPrefix("recond: ")

	var masterAddr = flag.String("masterAddr", "http://localhost:3000", "address of the recon-master server (along with protocol)")
	flag.Parse()

	conf, err := config.Init()
	if err != nil {
		log.Fatalln(err)
	}

	// agent represents a single agent on which the recond
	// is running.
	var agent = &Agent{
		UID: conf.UID,
	}

	err = agent.register(*masterAddr)
	if err != nil {
		log.Fatalln(err)
	}

	defer natsEncConn.Close()

	natsEncConn.Subscribe(agent.UID, func(s string) {
		fmt.Printf("Received a message: %s\n", s)
	})

	c := time.Tick(5 * time.Second)
	for now := range c {
		log.Println("Update sent at", now)
		if err := agent.update(*masterAddr); err != nil {
			log.Println(err)
		}
	}
}
