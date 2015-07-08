// Copyright 2015 Hari haran. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hariharan-uno/recon/blockdevice"
	"github.com/hariharan-uno/recon/cpu"
	"github.com/hariharan-uno/recon/languages"
	"github.com/hariharan-uno/recon/lsb"
	"github.com/hariharan-uno/recon/memory"
	"github.com/hariharan-uno/recon/uptime"
)

func main() {
	log.SetPrefix("recon: ")

	var addr = flag.String("addr", ":3030", "serve HTTP on `address`")
	flag.Parse()

	http.HandleFunc("/", reconHandler)
	fmt.Printf("recon: starting the server on http://localhost%s\n", *addr)
	fmt.Printf("recon: if you'd like the JSON to be indented, append %q to the above URL\n", "?indent=1")
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func copyMap(from, to map[string]interface{}) {
	for k, v := range from {
		to[k] = v
	}
}

// accumulateData accumulates data from all other packages.
// We just log the error but don't expose it, as we want
// to view memory data even if lsb data fails.
// TODO: Is it the right way? Rethink?
func accumulateData() map[string]interface{} {
	lsbdata, err := lsb.CollectData()
	if err != nil {
		log.Println(err)
	}
	memdata, err := memory.CollectData()
	if err != nil {
		log.Println(err)
	}
	cpudata, err := cpu.CollectData()
	if err != nil {
		log.Println(err)
	}
	blockdevicedata, err := blockdevice.CollectData()
	if err != nil {
		log.Println(err)
	}
	langsdata, err := languages.CollectData()
	if err != nil {
		log.Println(err)
	}
	uptimedata, err := uptime.CollectData()
	if err != nil {
		log.Println(err)
	}
	data := map[string]interface{}{
		"lsb":          lsbdata,
		"memory":       memdata,
		"cpu":          cpudata,
		"block_device": blockdevicedata,
		"languages":    langsdata,
		"recon_time":   time.Now(),
	}
	copyMap(uptimedata, data) // uptime Data is not namespaced.
	return data
}
