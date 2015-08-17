// Copyright 2015 CodeIgnition. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/codeignition/recon"
	"github.com/nats-io/nats"
)

// Agent is just recon.Agent. It has a separate type to
// add methods to it.
type Agent recon.Agent

func (a *Agent) register(addr string) error {
	if a.UID == "" {
		return errors.New("UID can't be empty")
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(a); err != nil {
		return err
	}

	// url.Parse instead of just appending will inform
	// about errors when addr or path is malformed.
	l, err := url.Parse(addr + agentsAPIPath)
	if err != nil {
		return err
	}
	resp, err := http.Post(l.String(), "application/json", &buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var t struct {
		NatsURL string `json:"nats_url"`
	}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&t); err != nil {
		return err
	}
	nc, err := nats.Connect(t.NatsURL)
	if err != nil {
		return err
	}
	// TODO: Should we return the conn instead of using a global?
	natsEncConn, err = nats.NewEncodedConn(nc, "json")
	if err != nil {
		return err
	}
	return nil
}

func (a *Agent) update(addr string) error {
	var buf bytes.Buffer

	m := recon.Metric{
		AgentUID: a.UID,
		Data:     accumulateData(),
	}

	if err := json.NewEncoder(&buf).Encode(&m); err != nil {
		return err
	}

	l, err := url.Parse(addr + metricsAPIPath)
	if err != nil {
		return err
	}
	resp, err := http.Post(l.String(), "application/json", &buf)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusCreated {
		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("response status code not %d; response body: %s\n", http.StatusCreated, b)
	}
	return nil
}
