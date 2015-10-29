// Copyright 2015 CodeIgnition. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

// +build linux

// Package disk provides disk metrics.
package disk

import (
	"os/exec"
	"strings"
)

type Data map[string]interface{}

func CollectData() (Data, error) {
	d := make(Data)
	d["disk"] = make(Data)
	disk := d["disk"].(Data)
	if err := sizeData(disk); err != nil {
		return d, err
	}
	return d, nil
}

func sizeData(d Data) error {
	out, err := exec.Command("df", "-P").Output()
	if err != nil {
		return err
	}

	lines := strings.Split(string(out), "\n")

	// lines[0] contains the column headings
	// Filesystem     1024-blocks     Used Available Capacity Mounted on
	for _, line := range lines[1:] {
		a := strings.Fields(line)
		if len(a) >= 6 {
			if strings.HasPrefix(a[0], "/dev/") {
				d[a[0]] = map[string]interface{}{
					"kb_size":         a[1],
					"kb_used":         a[2],
					"kb_available":    a[3],
					"percentage_used": a[4],
					"mounted_on":      a[5],
				}
			}
		}
	}
	return nil
}
