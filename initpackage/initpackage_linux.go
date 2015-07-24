// Copyright 2015 CodeIgnition. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

// Package initpackage exposes a single variable which
// says the init package the system is using, e.g. Systemd
package initpackage

import (
	"io/ioutil"
	"log"
	"strings"
)

// Name is the name of the init package used by the system
var Name string

func init() {
	if b, err := ioutil.ReadFile("/proc/1/comm"); err == nil {
		Name = strings.TrimSpace(string(b))
	} else {
		log.Println("initpackage: ", err)
	}
}
