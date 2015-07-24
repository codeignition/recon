// Copyright 2015 CodeIgnition. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

// Package blockdevice provides different block devices related data.
package blockdevice

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/codeignition/recon/internal/fileutil"
)

// Data represents the block devices data.
type Data map[string]interface{}

// CollectData collects the data and returns an error if any.
func CollectData() (Data, error) {
	d := make(Data)
	basePath := "/sys/block/"
	if !fileutil.Exists(basePath) {
		return nil, errors.New("blockdevice: cannot find /sys/block/ directory")
	}
	f, err := os.Open(basePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	devices, err := f.Readdirnames(-1)
	if err != nil {
		return nil, err
	}

	for _, device := range devices {
		d[device] = make(map[string]string)
		m := d[device].(map[string]string) // map for a device
		devicePath := filepath.Join(basePath, device)
		for _, key := range [...]string{"size", "removable"} {
			p := filepath.Join(devicePath, key)
			if fileutil.Exists(p) {
				b, err := ioutil.ReadFile(p)
				if err != nil {
					return nil, err
				}
				m[key] = strings.TrimSpace(string(b))
			}
		}

		for _, key := range [...]string{"model", "rev", "state", "timeout", "vendor"} {
			p := filepath.Join(devicePath, "device", key)
			if fileutil.Exists(p) {
				b, err := ioutil.ReadFile(p)
				if err != nil {
					return nil, err
				}
				m[key] = strings.TrimSpace(string(b))
			}
		}

		p := filepath.Join(devicePath, "queue", "rotational")
		if fileutil.Exists(p) {
			b, err := ioutil.ReadFile(p)
			if err != nil {
				return nil, err
			}
			m["rotational"] = strings.TrimSpace(string(b))
		}
	}
	return d, nil
}
