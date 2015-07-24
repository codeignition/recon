// Copyright 2015 CodeIgnition. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

// Package cpu gives various CPU stats.
package cpu

import (
	"bufio"
	"os"
	"strings"
)

// Data represents the CPU data.
type Data map[string]interface{}

// CollectData collects the data and returns an error if any.
func CollectData() (Data, error) {
	d := make(Data)
	f, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	var (
		p     string                 // current processor
		pm    map[string]interface{} // current processor map
		total int                    // total number of processors

	)
	real := make(map[string]struct{}) // real number of processors.
	// It is a map so that physical ids are keys for existence check.
	// We use a struct{} instead of bool to not allocate memory.

	for s.Scan() {
		l := strings.Split(s.Text(), ":")
		if len(l) == 2 {
			k, v := strings.TrimSpace(l[0]), strings.TrimSpace(l[1])
			switch k {
			case "processor":
				p = v
				d[p] = make(map[string]interface{})
				pm = d[p].(map[string]interface{})
				total++
			case "vendor_id":
				pm["vendor_id"] = v
			case "cpu family":
				pm["family"] = v
			case "model":
				pm["model"] = v
			case "model name":
				pm["model_name"] = v
			case "stepping":
				pm["stepping"] = v
			case "cpu MHz":
				pm["mhz"] = v
			case "cache size":
				pm["cache_size"] = v
			case "physical id":
				pm["physical_id"] = v
				real[v] = struct{}{}
			case "core id":
				pm["core_id"] = v
			case "cpu cores":
				pm["cores"] = v
			case "flags":
				pm["flags"] = strings.Split(v, " ")

			}
		}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	d["total"] = total
	d["real"] = len(real)
	return d, nil
}
