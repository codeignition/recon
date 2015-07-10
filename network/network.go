// Copyright 2015 Hari haran. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package network

import (
	"os/exec"
	"strings"
)

// Data represents the network data.
type Data map[string]interface{}

// CollectData collects the data and returns an error if any.
func CollectData() (Data, error) {
	d := make(Data)
	if out, err := exec.Command("route", "-n").Output(); err == nil {
		s := strings.Split(string(out), "\n")
		// s[0] is the title, s[1] is the column headings. Also, we only
		// consider s[2] for the default interface and gateway.
		a := strings.Fields(s[2])
		d["default_gateway"] = a[1]
		d["default_interface"] = a[7]
	} else {
		return nil, err
	}
	return d, nil
}
