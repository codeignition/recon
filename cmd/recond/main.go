// Copyright 2015 CodeIgnition. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"flag"
	"log"
	"sync"

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

// ctxCancelFunc stores the map of policy name to
// the context cancel function.
var ctxCancelFunc = struct {
	sync.Mutex
	m map[string]context.CancelFunc
}{
	m: make(map[string]context.CancelFunc),
}

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

	natsEncConn.Subscribe(agent.UID+"_add_policy", AddPolicyHandler(conf))
	natsEncConn.Subscribe(agent.UID+"_delete_policy", DeletePolicyHandler(conf))
	natsEncConn.Subscribe(agent.UID+"_modify_policy", ModifyPolicyHandler(conf))

	// this is just to block the main function from exiting
	c := make(chan struct{})
	<-c
}

func runStoredPolicies(c *config.Config) {
	for _, p := range c.PolicyConfig {
		log.Printf("adding the policy %s...", p.Name)
		go func(p policy.Policy) {
			ctx, cancel := context.WithCancel(context.Background())
			events, err := p.Execute(ctx)
			if err != nil {
				log.Print(err) // TODO: send to a nats errors channel
			}
			ctxCancelFunc.Lock()
			ctxCancelFunc.m[p.Name] = cancel
			ctxCancelFunc.Unlock()

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

func deletePolicy(c *config.Config, policyName string) error {
	defer ctxCancelFunc.Unlock()
	ctxCancelFunc.Lock()

	if _, ok := ctxCancelFunc.m[policyName]; !ok {
		return errors.New("policy not found")
	}
	log.Printf("deleting the policy %s...", policyName)

	delete(ctxCancelFunc.m, policyName)
	defer c.Unlock()
	c.Lock()
	for i, q := range c.PolicyConfig {
		if q.Name == policyName {
			c.PolicyConfig = append(c.PolicyConfig[:i], c.PolicyConfig[i+1:]...)
		}
	}

	if err := c.Save(); err != nil {
		return err
	}
	return nil
}
