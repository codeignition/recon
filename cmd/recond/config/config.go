// Copyright 2015 CodeIgnition. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

// Package config provides utilities for dealing
// with configuration files for recond.
package config

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"sync"

	"github.com/codeignition/recon/internal/fileutil"
	"github.com/codeignition/recon/policy"
)

const configFileName = ".recond.json"

// config file path in the local machine
var configPath string

func init() {
	usr, err := user.Current()
	if err != nil {
		log.Fatalln(err)
	}
	configPath = filepath.Join(usr.HomeDir, configFileName)
}

// Config represents the configuration for the recond
// running on a particular machine.
type Config struct {
	sync.Mutex
	UID          string // Unique Identifier to register with marksman
	PolicyConfig policy.Config
}

// Init initializes and returns a Config. i.e. if the config file doesn't exist,
// it generates an new UID, creates the config file and returns the corresponding Config.
func Init() (*Config, error) {
	if fileutil.Exists(configPath) {
		f, err := os.Open(configPath)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		return parseConfig(f)
	}
	uid, err := generateUID()
	if err != nil {
		return nil, err
	}
	c := &Config{
		UID: uid,
	}
	err = c.Save()
	return c, err
}

// Save saves the config in a configuration file.
// If it already exists, it removes it and writes it freshly.
// Make sure you lock and unlock the config while calling Save.
func (c *Config) Save() error {
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

	var out bytes.Buffer
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	out.Write(b)
	out.WriteTo(f)

	if err := f.Close(); err != nil {
		return err
	}

	return nil
}

func (c *Config) AddPolicy(p policy.Policy) error {
	defer c.Unlock()
	c.Lock()
	for _, k := range c.PolicyConfig {
		if k.Name == p.Name {
			return errors.New("policy with the given name already exists")
		}
	}
	c.PolicyConfig = append(c.PolicyConfig, p)
	return nil
}

// parseConfig reads from a io.Reader and
// creates a Config struct accordingly.
// It takes an io.Reader so that it is easier
// to test it.
func parseConfig(r io.Reader) (*Config, error) {
	dec := json.NewDecoder(r)
	var c Config
	if err := dec.Decode(&c); err != nil {
		return nil, err
	}
	return &c, nil
}

// generateUID returns a UID string of length 12.
func generateUID() (string, error) {
	b := make([]byte, 6)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	// b is of 6 bytes length, but when converted into base-16,
	// uid has a length of 12
	uid := fmt.Sprintf("%x", b)
	return uid, nil
}
