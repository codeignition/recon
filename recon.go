// Copyright 2015 Hari haran. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"log"
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

	b, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		fmt.Println("error:", err)
	}

	fmt.Printf("%s\n", b)
}

func copyMap(from, to map[string]interface{}) {
	for k, v := range from {
		to[k] = v
	}
}
