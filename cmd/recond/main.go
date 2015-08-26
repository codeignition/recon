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
	"github.com/codeignition/recon/policy"
	"github.com/nats-io/nats"
)

const agentsAPIPath = "/api/agents" // agents path in the marksman server

// natsEncConn is the opened with the URL obtained from marksman.
// It is populated if the agent registers successfully.
var natsEncConn *nats.EncodedConn

// updateInterval is time.Duration that specifies the interval
// between two consecutive updates.
const updateInterval = 5 * time.Second

func main() {
	log.SetPrefix("recond: ")

	var marksmanAddr = flag.String("marksman", "http://localhost:3000", "address of the marksman server")
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

	err = agent.register(*marksmanAddr)
	if err != nil {
		log.Fatalln(err)
	}

	defer natsEncConn.Close()

	go runStoredPolicies(conf)

	natsEncConn.Subscribe(agent.UID, func(s string) {
		fmt.Printf("Received a message: %s\n", s)
	})

	natsEncConn.Subscribe(agent.UID+"_policy", func(subj, reply string, p *policy.Policy) {
		fmt.Printf("Received a Policy: %v\n", p)
		if err := conf.AddPolicy(*p); err != nil {
			natsEncConn.Publish(reply, err.Error())
			return
		}
		if err := conf.Save(); err != nil {
			natsEncConn.Publish(reply, err.Error())
			return
		}
		if err := p.Execute(); err != nil {
			natsEncConn.Publish(reply, err.Error())
			return
		}
		natsEncConn.Publish(reply, "policy ack") // acknowledge
	})

	c := time.Tick(updateInterval)
	for now := range c {
		log.Println("Update sent at", now)
		agent.update()
	}
}

func runStoredPolicies(c *config.Config) {
	for _, p := range c.PolicyConfig {
		go func() {
			if err := p.Execute(); err != nil {
				log.Fatal(err)
			}
		}()
	}
}
