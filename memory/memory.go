package memory

import (
	"bufio"
	"os"
	"strings"
)

type Data map[string]interface{}

func CollectData() (Data, error) {
	d := make(Data)
	d["swap"] = make(map[string]string)
	swap, _ := d["swap"].(map[string]string)
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return d, err
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		l := strings.Split(s.Text(), ":")
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
		case "Sunreclaim":
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
	if err := s.Err(); err != nil {
		return d, err
	}
	return d, nil
}
