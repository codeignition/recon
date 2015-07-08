// Copyright 2015 Hari haran. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

// Package languages gives various languages related data.
package languages

import (
	"os/exec"
	"strings"
)

// Data represents the languages data.
type Data map[string]interface{}

// CollectData collects the data and returns an error if any.
func CollectData() (Data, error) {
	d := make(Data)
	goData(d)
	perlData(d)
	pythonData(d)
	// TODO: ruby, c ?
	return d, nil
}

func goData(d map[string]interface{}) {
	if out, err := exec.Command("go", "version").Output(); err == nil {
		lines := strings.Split(string(out), " ")
		d["go"] = make(map[string]string)
		m := d["go"].(map[string]string)
		m["version"] = lines[2][2:]
	}
}

func perlData(d map[string]interface{}) {
	if out, err := exec.Command("perl", "-V:version", "-V:archname").Output(); err == nil {
		lines := strings.Split(string(out), "\n")
		d["perl"] = make(map[string]string)
		m := d["perl"].(map[string]string)
		for _, line := range lines {
			l := strings.Split(line, "=")
			if len(l) == 2 {
				v := strings.Trim(l[1], `';`)
				switch l[0] {
				case "version":
					m["version"] = v
				case "archname":
					m["archname"] = v
				}
			}
		}
	}
}

func pythonData(d map[string]interface{}) {
	if out, err := exec.Command("python", "-c", "import sys; print(sys.version)").Output(); err == nil {
		// only the first line is required.
		line := strings.Split(string(out), "\n")[0] // length check necessary?
		d["python"] = make(map[string]string)
		m := d["python"].(map[string]string)
		l := strings.SplitN(line, " ", 2)
		if len(l) == 2 {
			m["version"] = l[0]
			m["builddate"] = strings.Trim(l[1], "(default, ) ")
		}
	}
}
