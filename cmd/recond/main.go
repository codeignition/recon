// Copyright 2015 CodeIgnition. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"

	"golang.org/x/net/context"

	"github.com/codeignition/recon/cmd/recond/config"
	"github.com/codeignition/recon/policy"
	_ "github.com/codeignition/recon/policy/handlers"
	"github.com/nats-io/nats"
)

const agentsAPIPath = "/api/agents" // agents path in the marksman server

// natsEncConn is the opened with the URL obtained from marksman.
// It is populated if the agent registers successfully.
var natsEncConn *nats.EncodedConn

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

	if err := addSystemDataPolicy(conf); err != nil {
		log.Fatal(err)
	}

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
		events, err := p.Execute(context.TODO())
		if err != nil {
			natsEncConn.Publish(reply, err.Error())
			return
		}
		natsEncConn.Publish(reply, "policy ack") // acknowledge
		for e := range events {
			natsEncConn.Publish("policy_events", e)
		}
	})

	// this is just to block the main function from exiting
	c := make(chan struct{})
	<-c
}

func runStoredPolicies(c *config.Config) {
	for _, p := range c.PolicyConfig {
		log.Printf("adding the policy %s...", p.Name)
		go func(p policy.Policy) {
			events, err := p.Execute(context.TODO())
			if err != nil {
				log.Print(err) // TODO: send to a nats errors channel
			}
			for e := range events {
				natsEncConn.Publish("policy_events", e)
			}
		}(p)
	}
}

func addSystemDataPolicy(c *config.Config) error {
	// if the policy already exists, return silently
	for _, p := range c.PolicyConfig {
		if p.Name == "default_system_data" {
			return nil
		}
	}

	p := policy.Policy{
		Name:     "default_system_data",
		AgentUID: c.UID,
		Type:     "system_data",
		M: map[string]string{
			"interval": "5s",
		},
	}
	if err := c.AddPolicy(p); err != nil {
		return err

	}
	if err := c.Save(); err != nil {
		return err
	}
	return nil
}
