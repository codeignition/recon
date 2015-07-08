// Copyright 2015 Hari haran. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

// Package uptime gives uptime and idletime data.
package uptime

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
)

// Data holds uptime and idletime data.
type Data map[string]interface{}

// CollectData collects the data and returns an error if any.
func CollectData() (Data, error) {
	d := make(Data)
	b, err := ioutil.ReadFile("/proc/uptime")
	if err != nil {
		log.Println("uptime: ", err)
		return nil, err
	} else {
		s := strings.Split(strings.Trim(string(b), "\n"), " ")
		if len(s) == 2 {
			us, err := strconv.ParseFloat(s[0], 32)
			if err != nil {
				return nil, err
			}
			d["uptime_seconds"] = us
			dur := time.Duration(us) * time.Second
			d["uptime"] = dur.String()

			is, err := strconv.ParseFloat(s[1], 32)
			if err != nil {
				return nil, err
			}
			d["idletime_seconds"] = is
			dur = time.Duration(is) * time.Second
			d["idletime"] = dur.String()
		}
	}
	return d, nil
}
