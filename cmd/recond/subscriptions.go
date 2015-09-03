// Copyright 2015 CodeIgnition. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package main

import (
	"log"

	"github.com/codeignition/recon/cmd/recond/config"
	"github.com/codeignition/recon/policy"
	"golang.org/x/net/context"
)

func AddPolicyHandler(conf *config.Config) func(subj, reply string, p *policy.Policy) {
	return func(subj, reply string, p *policy.Policy) {
		log.Printf("policy_add received: %s\n", p.Name)
		if err := conf.AddPolicy(*p); err != nil {
			natsEncConn.Publish(reply, err.Error())
			return
		}
		if err := conf.Save(); err != nil {
			natsEncConn.Publish(reply, err.Error())
			return
		}
		ctx, cancel := context.WithCancel(context.Background())
		events, err := p.Execute(ctx)
		if err != nil {
			natsEncConn.Publish(reply, err.Error())
			return
		}
		ctxCancelFunc.Lock()
		ctxCancelFunc.m[p.Name] = cancel
		ctxCancelFunc.Unlock()

		natsEncConn.Publish(reply, "policy_add_ack") // acknowledge policy add
		for e := range events {
			natsEncConn.Publish("policy_events", e)
		}
	}
}

func DeletePolicyHandler(conf *config.Config) func(subj, reply string, p *policy.Policy) {
	return func(subj, reply string, p *policy.Policy) {
		log.Printf("policy_delete received: %s\n", p.Name)
		ctxCancelFunc.Lock()
		cancel := ctxCancelFunc.m[p.Name]
		ctxCancelFunc.Unlock()
		cancel()
		if err := deletePolicy(conf, p.Name); err != nil {
			natsEncConn.Publish(reply, err.Error())
			return
		}
		natsEncConn.Publish(reply, "policy_delete_ack") // acknowledge policy delete
	}
}

func ModifyPolicyHandler(conf *config.Config) func(subj, reply string, p *policy.Policy) {
	return func(subj, reply string, p *policy.Policy) {
		log.Printf("modify_policy received: %s\n", p.Name)

		// We receive the complete policy with the new values
		// and delete the old policy and stop its execution.
		// Then we add the new policy.
		ctxCancelFunc.Lock()
		cancel := ctxCancelFunc.m[p.Name]
		ctxCancelFunc.Unlock()
		cancel()
		if err := deletePolicy(conf, p.Name); err != nil {
			log.Print(err)
			natsEncConn.Publish(reply, err.Error())
			return
		}
		log.Printf("adding the policy %s...", p.Name)
		if err := conf.AddPolicy(*p); err != nil {
			natsEncConn.Publish(reply, err.Error())
			return
		}
		if err := conf.Save(); err != nil {
			natsEncConn.Publish(reply, err.Error())
			return
		}
		ctx, cancel := context.WithCancel(context.Background())
		events, err := p.Execute(ctx)
		if err != nil {
			natsEncConn.Publish(reply, err.Error())
			return
		}
		ctxCancelFunc.Lock()
		ctxCancelFunc.m[p.Name] = cancel
		ctxCancelFunc.Unlock()

		natsEncConn.Publish(reply, "modify_policy_ack") // acknowledge policy delete
		for e := range events {
			natsEncConn.Publish("policy_events", e)
		}
	}
}
