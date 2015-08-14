// Copyright 2015 CodeIgnition. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/codeignition/recon"
	"github.com/codeignition/recon/internal/fileutil"
	"github.com/nats-io/nats"
)

const (
	metricsAPIPath = "/api/metrics" // metrics path in the master server
	agentsAPIPath  = "/api/agents"  // agents path in the master server
)

// natsEncConn is the opened with the URL obtained from marksman.
// It is populated if the agent registers successfully.
var natsEncConn *nats.EncodedConn

// config file path in the local machine
var configPath string

// Agent is just recon.Agent. It has a separate type to
// add methods to it.
type Agent recon.Agent

// Config represents the configuration for the recond
// running on a particular machine.
type Config struct {
	UID string `json:"uid"` // Unique Identifier to register with marksman
}

func init() {
	usr, err := user.Current()
	if err != nil {
		log.Fatalln(err)
	}
	configPath = filepath.Join(usr.HomeDir, ".recond.json")
}

// save saves the config in the configPath
// If it already exists, it removes it and writes it freshly.
func (c *Config) save() error {
	if fileutil.Exists(configPath) {
		if err := os.Remove(configPath); err != nil {
			return err
		}
	}

	f, err := os.Create(configPath)
	if err != nil {
		return err
	}
	// Brad Fitzpatrick:
	// Never defer a file Close when the file was opened for writing.
	// Many filesystems do their real work (and thus their real failures) on close.
	// You can defer a file.Close for Read, but not for write.

	enc := json.NewEncoder(f)
	if err := enc.Encode(c); err != nil {
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}

	return nil
}

func parseConfigFile(file string) (*Config, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	var conf Config
	if err := dec.Decode(&conf); err != nil {
		return nil, err
	}
	return &conf, nil
}

// initConfig returns a Config. If the config file doesn't exist,
// it creates it and returns the corresponding Config.
func initConfig() (*Config, error) {
	if fileutil.Exists(configPath) {
		return parseConfigFile(configPath)
	}

	uid, err := generateUID()
	if err != nil {
		return nil, err
	}
	conf := &Config{
		UID: uid,
	}
	err = conf.save()
	return conf, err
}

func main() {
	log.SetPrefix("recond: ")

	var masterAddr = flag.String("masterAddr", "http://localhost:3000", "address of the recon-master server (along with protocol)")
	flag.Parse()

	conf, err := initConfig()
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
		NatsUrl string `json:"nats_url"`
	}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&t); err != nil {
		return err
	}
	nc, err := nats.Connect(t.NatsUrl)
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

func generateUID() (string, error) {
	b := make([]byte, 6)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	uid := fmt.Sprintf("%x", b)
	return uid, nil
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
