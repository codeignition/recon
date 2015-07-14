// Copyright 2015 Hari haran. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

// +build linux

package filesystem

import (
	"bufio"
	"os"
	"os/exec"
	"strings"

	"github.com/hariharan-uno/recon/internal/fileutil"
)

// Data represents the filesystem data.
type Data map[string]interface{}

// CollectData collects the data and returns an error if any.
func CollectData() (Data, error) {
	d := make(Data)
	if err := sizeData(d); err != nil {
		return d, err
	}
	if err := inodeData(d); err != nil {
		return d, err
	}
	if err := mountData(d); err != nil {
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
			d[a[0]] = map[string]interface{}{
				"kb_size":         a[1],
				"kb_used":         a[2],
				"kb_available":    a[3],
				"percentage_used": a[4],
				"mounted_on":      a[5],
			}
		}
	}
	return nil
}

func inodeData(d Data) error {
	out, err := exec.Command("df", "-iP").Output()
	if err != nil {
		return err
	}

	lines := strings.Split(string(out), "\n")

	// lines[0] contains the column headings
	// Filesystem       Inodes  IUsed    IFree IUse% Mounted on
	for _, line := range lines[1:] {
		a := strings.Fields(line)
		if len(a) >= 6 {
			m := d[a[0]].(map[string]interface{})
			m["total_inodes"] = a[1]
			m["inodes_used"] = a[2]
			m["inodes_available"] = a[3]
			m["inodes_percentage_used"] = a[4]
			m["mount"] = a[5]
		}
	}
	return nil
}

func mountData(d Data) error {
	out, err := exec.Command("mount").Output()
	if err != nil {
		return err
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		// proc on /proc type proc (rw,noexec,nosuid,nodev)

		a := strings.Fields(line)
		if len(a) >= 6 {
			if _, ok := d[a[0]]; !ok {
				d[a[0]] = make(map[string]interface{})
			}
			m := d[a[0]].(map[string]interface{})
			m["mount"] = a[2]
			m["fs_type"] = a[4]
			m["mount_options"] = strings.Split(strings.Trim(a[5], "()"), ",")
		}
	}

	// Get missing mount info from /proc/mounts

	if fileutil.Exists("/proc/mounts") {
		f, err := os.Open("/proc/mounts")
		if err != nil {
			return err
		}
		defer f.Close()
		s := bufio.NewScanner(f)
		for s.Scan() {
			a := strings.Fields(s.Text())
			if len(a) >= 4 {
				if _, ok := d[a[0]]; !ok {
					d[a[0]] = map[string]interface{}{
						"mount":         a[1],
						"fs_type":       a[2],
						"mount_options": strings.Split(a[3], ","),
					}
				}
			}
		}
		if err := s.Err(); err != nil {
			return err
		}
	}
	return nil
}
