// Copyright 2015 Hari haran. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

// Package memory provides different memory related data.
package memory

import (
	"bufio"
	"os"
	"strings"
)

// Data represents the memory data.
type Data map[string]interface{}

// CollectData collects the data and returns an error if any.
func CollectData() (Data, error) {

	// allocate memory
	d := make(Data)
	d["swap"] = make(map[string]string)

	// cast interface{} to a map
	swap := d["swap"].(map[string]string)

	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return d, err
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		l := strings.Split(s.Text(), ":")
		if len(l) == 2 {
			k, v := l[0], strings.TrimSpace(l[1])
			switch k {
			case "MemTotal":
				d["total"] = v
			case "MemFree":
				d["free"] = v
			case "MemAvailable":
				d["available"] = v
			case "Buffers":
				d["buffers"] = v
			case "Cached":
				d["cached"] = v
			case "SwapCached":
				swap["cached"] = v
			case "Active":
				d["active"] = v
			case "Inactive":
				d["inactive"] = v
			case "SwapTotal":
				swap["total"] = v
			case "SwapFree":
				swap["free"] = v
			case "Dirty":
				d["dirty"] = v
			case "Writeback":
				d["writeback"] = v
			case "AnonPages":
				d["anon_pages"] = v
			case "Mapped":
				d["mapped"] = v
			case "Slab":
				d["slab"] = v
			case "SReclaimable":
				d["slab_reclaimable"] = v
			case "SUnreclaim":
				d["slab_unreclaim"] = v
			case "PageTables":
				d["page_tables"] = v
			case "NFS_Unstable":
				d["nfs_unstable"] = v
			case "Bounce":
				d["bounce"] = v
			case "CommitLimit":
				d["commit_limit"] = v
			case "Committed_AS":
				d["committed_as"] = v
			case "VmallocTotal":
				d["vmalloc_total"] = v
			case "VmallocUsed":
				d["vmalloc_used"] = v
			case "VmallocChunk":
				d["vmalloc_chunk"] = v
			}
		}
	}
	if err := s.Err(); err != nil {
		return d, err
	}
	return d, nil
}
