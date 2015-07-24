// Copyright 2015 CodeIgnition. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

// +build linux

package counters

import (
	"os/exec"
	"strings"
)

// Data represents the counters data.
type Data map[string]interface{}

// CollectData collects the data and returns an error if any.
func CollectData() (Data, error) {
	d := make(Data)

	d["network"] = make(map[string]interface{})
	network := d["network"].(map[string]interface{})

	network["interfaces"] = make(map[string]interface{})
	ifaces := network["interfaces"].(map[string]interface{})

	out, err := exec.Command("ip", "-d", "-s", "link").Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(out), "\n")
	for i := 0; i < len(lines); {
		if strings.ContainsAny(lines[i], "<>") {
			rx := strings.Fields(lines[i+3]) // bytes  packets  errors  dropped overrun mcast
			tx := strings.Fields(lines[i+5]) // bytes  packets  errors  dropped carrier collsns

			k := strings.TrimSpace(strings.Split(lines[i], ":")[1]) // [0] is the index
			ifaces[k] = map[string]map[string]string{
				"rx": {
					"bytes":   rx[0],
					"packets": rx[1],
					"errors":  rx[2],
					"dropped": rx[3],
					"overrun": rx[4],
					"mcast":   rx[5],
				},
				"tx": {
					"bytes":      tx[0],
					"packets":    tx[1],
					"errors":     tx[2],
					"dropped":    tx[3],
					"carrier":    tx[4],
					"collisions": tx[5],
				},
			}

			i += 6
		} else {
			i++
		}
	}
	return d, nil
}
