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
type Data struct {
	ID,
	Release,
	Codename,
	Description string
}

// CollectData collects the data and returns an error if any.
func CollectData() (*Data, error) {
	d := &Data{}

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
				d.ID = v
			case "DISTRIB_RELEASE":
				d.Release = v
			case "DISTRIB_CODENAME":
				d.Codename = v
			case "DISTRIB_DESCRIPTION":
				d.Description = strings.Trim(v, `"`)
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
					d.ID = v
				case "Release":
					d.Release = v
				case "Codename":
					d.Codename = v
				case "Description":
					d.Description = v
				}
			}
		}
		return d, nil
	}

	return nil, errors.New("cannot find /etc/lsb-release or /usr/bin/lsb_release")
}
