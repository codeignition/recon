// Copyright 2015 CodeIgnition. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package etc

import (
	"bufio"
	"os"
	"strings"
)

// Data represents the etc data.
type Data map[string]interface{}

// CollectData collects the data and returns an error if any.
func CollectData() (Data, error) {
	d := make(Data)

	d["passwd"] = make(map[string]interface{})
	passwd := d["passwd"].(map[string]interface{})
	f, err := os.Open("/etc/passwd")
	if err != nil {
		return d, err
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		l := strings.Split(line, ":")
		if len(l) == 7 {
			passwd[l[0]] = map[string]string{
				"uid":   l[2],
				"gid":   l[3],
				"gecos": l[4],
				"dir":   l[5],
				"shell": l[6],
			}
		}
	}
	if err := s.Err(); err != nil {
		return d, err
	}

	d["group"] = make(map[string]interface{})
	group := d["group"].(map[string]interface{})
	f, err = os.Open("/etc/group")
	if err != nil {
		return d, err
	}
	defer f.Close()
	s = bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		l := strings.Split(line, ":")
		if len(l) == 4 {
			group[l[0]] = map[string]interface{}{
				"gid":     l[2],
				"members": strings.Split(l[3], ","),
			}
		}
	}
	if err := s.Err(); err != nil {
		return d, err
	}
	return d, nil
}
