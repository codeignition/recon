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
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/codeignition/recon/internal/fileutil"
)

const (
	metricsPath = "/metrics.json" // metrics path in the master server
	agentsPath  = "/agents.json"  // agents path in the master server
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

	if err := registerAgent(*masterAddr, uid); err != nil {
		log.Fatalln(err)
	}

	c := time.Tick(5 * time.Second)
	for now := range c {
		log.Println("Update sent at", now)
		if err := update(*masterAddr); err != nil {
			log.Println(err)
		}
	}
}

func registerAgent(addr, uid string) error {
	var buf bytes.Buffer
	d := map[string]interface{}{
		"agent": map[string]string{
			"uid": uid,
		},
	}
	if err := json.NewEncoder(&buf).Encode(&d); err != nil {
		return err
	}
	resp, err := http.Post(addr+agentsPath, "application/json", &buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// TODO: Don't print the response, but store the messaging server URL and subscribe to it.
	fmt.Printf("%s\n", resp.Body)
	return nil
}

func generateUID() (string, error) {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	uid := fmt.Sprintf("%x", b)
	return uid, nil
}

func update(addr string) error {
	var buf bytes.Buffer
	d := map[string]interface{}{
		"metric": accumulateData(),
	}
	if err := json.NewEncoder(&buf).Encode(&d); err != nil {
		return err
	}
	resp, err := http.Post(addr+metricsPath, "application/json", &buf)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("response status code not %d; response is %v\n", http.StatusCreated, resp)
	}
	return nil
}
