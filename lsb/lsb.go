// Copyright 2015 Hari haran. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

// +build linux

// lsb is used to identify the distribution being used
// and its compliance with Linux Standard Base.
package lsb

import (
	"bufio"
	"errors"
	"os"
	"os/exec"
	"strings"

	"github.com/hariharan-uno/recon/internal/fileutil"
)

// Data represents the lsb data.
type Data map[string]string

// CollectData collects the data and returns an error if any.
func CollectData() (Data, error) {
	d := make(Data)
	if fileutil.Exists("/etc/lsb-release") {
		f, err := os.Open("/etc/lsb-release")
		if err != nil {
			return d, err
		}
		defer f.Close()
		s := bufio.NewScanner(f)
		for s.Scan() {
			l := strings.Split(s.Text(), "=")
			k, v := l[0], l[1]
			switch k {
			case "DISTRIB_ID":
				d["id"] = v
			case "DISTRIB_RELEASE":
				d["release"] = v
			case "DISTRIB_CODENAME":
				d["codename"] = v
			case "DISTRIB_DESCRIPTION":
				d["description"] = strings.Trim(v, `"`)
			}
		}
		if err := s.Err(); err != nil {
			return d, err
		}
		return d, nil
	}

	if fileutil.Exists("/usr/bin/lsb_release") {
		out, err := exec.Command("/usr/bin/lsb_release", "-a").Output()
		if err != nil {
			return d, err
		}
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			l := strings.Split(line, ":")
			if len(l) == 2 {
				k, v := l[0], strings.TrimSpace(l[1])
				switch k {
				case "Distributor ID":
					d["id"] = v
				case "Release":
					d["release"] = v
				case "Codename":
					d["codename"] = v
				case "Description":
					d["description"] = v
				}
			}
		}
		return d, nil
	}

	return nil, errors.New("cannot find /etc/lsb-release or /usr/bin/lsb_release")
}
