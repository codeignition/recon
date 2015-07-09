// Copyright 2015 Hari haran. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

// Package kernel provides different kernel related data.
package kernel

import (
	"os/exec"
	"strings"
)

// Data represents the kernel data.
type Data map[string]interface{}

// unameArgs maps the corresponding
// argument for the uname command.
var unameArgs = map[string]string{
	"name":    "-s",
	"release": "-r",
	"version": "-v",
	"machine": "-m",
	"os":      "-o",
}

// CollectData collects the data and returns an error if any.
func CollectData() (Data, error) {
	d := make(Data)
	for k, v := range unameArgs {
		out, err := exec.Command("uname", v).Output()
		if err != nil {
			return nil, err
		}
		s := strings.TrimSpace(string(out))
		d[k] = s
		out, err = exec.Command("env", "lsmod").Output()
		if err != nil {
			return nil, err
		}
		lines := strings.Split(string(out), "\n")

		// lines[0] contains column headings.
		for _, line := range lines[1:] {
			l := strings.Fields(line)
			if len(l) >= 3 {
				d[l[0]] = map[string]string{
					"size":     l[1],
					"refcount": l[2],
				}
			}
		}

	}
	return d, nil
}
