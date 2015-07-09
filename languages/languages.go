// Copyright 2015 Hari haran. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

// Package languages gives various languages related data.
package languages

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/hariharan-uno/recon/internal/fileutil"
)

// Data represents the languages data.
type Data map[string]interface{}

// CollectData collects the data and returns an error if any.
func CollectData() (Data, error) {
	d := make(Data)
	goData(d)
	perlData(d)
	pythonData(d)
	rubyData(d)
	// TODO: c ?
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

func rubyData(d map[string]interface{}) {
	rvars := map[string]string{
		"platform":      "RUBY_PLATFORM",
		"version":       "RUBY_VERSION",
		"release_date":  "RUBY_RELEASE_DATE",
		"target":        "RbConfig::CONFIG['target']",
		"target_cpu":    "RbConfig::CONFIG['target_cpu']",
		"target_vendor": "RbConfig::CONFIG['target_vendor']",
		"target_os":     "RbConfig::CONFIG['target_os']",
		"host":          "RbConfig::CONFIG['host']",
		"host_cpu":      "RbConfig::CONFIG['host_cpu']",
		"host_os":       "RbConfig::CONFIG['host_os']",
		"host_vendor":   "RbConfig::CONFIG['host_vendor']",
		"bin_dir":       "RbConfig::CONFIG['bindir']",
		"ruby_bin":      "::File.join(RbConfig::CONFIG['bindir'], RbConfig::CONFIG['ruby_install_name'])",
	}
	d["ruby"] = make(map[string]string)
	m := d["ruby"].(map[string]string)
	for k, v := range rvars {
		t := fmt.Sprintf(`require "rbconfig"; puts %s`, v)
		if out, err := exec.Command("ruby", "-e", t).Output(); err == nil {
			m[k] = strings.TrimSpace(string(out))
		}
	}
	if out, err := exec.Command("ruby", "-e", `require "rubygems"; puts Gem::default_exec_format % "gem"`).Output(); err == nil {
		g := strings.TrimSpace(string(out))
		var gemBin string
		if p := filepath.Join(m["bin_dir"], g); fileutil.Exists(p) {
			gemBin = p
		} else if p := filepath.Join(m["bin_dir"], "gem"); fileutil.Exists(p) {
			gemBin = p
		}
		if gemBin != "" {
			m["gem_bin"] = gemBin

			// TODO: A bit of a doubt. Check gems_dir once.
			if out, err := exec.Command(m["ruby_bin"], gemBin, "env", "gemdir").Output(); err == nil {
				m["gems_dir"] = strings.TrimSpace(string(out))
			}
		}
	}
}
