// Copyright 2015 CodeIgnition. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
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
)

const (
	metricsAPIPath = "/api/metrics" // metrics path in the master server
	agentsAPIPath  = "/api/agents"  // agents path in the master server
)

// config file path in the local machine
var configPath string

func init() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	configPath = filepath.Join(usr.HomeDir, ".recon")
	if fileutil.Exists(configPath) {
		// TODO: Here we are deleting the file, while development.
		// Change the logic to get the uid from this file later.
		err := os.Remove(configPath)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func main() {
	log.SetPrefix("recon: ")

	var masterAddr = flag.String("masterAddr", "http://localhost:3000", "address of the recon-master server (along with protocol)")
	flag.Parse()

	uid, err := generateUID()
	if err != nil {
		log.Fatalln(err)
	}

	f, err := os.Create(configPath)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	_, err = f.WriteString("uid: " + uid + "\n")
	if err != nil {
		log.Fatalln(err)
	}

	a, err := registerAgent(*masterAddr, uid)
	if err != nil {
		log.Fatalln(err)
	}

	c := time.Tick(5 * time.Second)
	for now := range c {
		log.Println("Update sent at", now)
		if err := update(*masterAddr, a); err != nil {
			log.Println(err)
		}
	}
}

func registerAgent(addr, uid string) (*recon.Agent, error) {
	var buf bytes.Buffer
	a := &recon.Agent{UID: uid}
	if err := json.NewEncoder(&buf).Encode(a); err != nil {
		return nil, err
	}

	// url.Parse instead of just appending will inform
	// about errors if the code changes and the url is malformed.
	l, err := url.Parse(addr + agentsAPIPath)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(l.String(), "application/json", &buf)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// TODO: Don't print the response, but store the messaging server URL and subscribe to it.
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%s\n", string(contents))
	return a, nil
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

func update(addr string, a *recon.Agent) error {
	var buf bytes.Buffer
	d := accumulateData()
	d["recon_uid"] = a.UID
	if err := json.NewEncoder(&buf).Encode(&d); err != nil {
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
		return fmt.Errorf("response status code not %d; response is %v\n", http.StatusCreated, resp)
	}
	return nil
}
