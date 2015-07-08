// Copyright 2015 Hari haran. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hariharan-uno/recon/blockdevice"
	"github.com/hariharan-uno/recon/cpu"
	"github.com/hariharan-uno/recon/languages"
	"github.com/hariharan-uno/recon/lsb"
	"github.com/hariharan-uno/recon/memory"
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

	// The field names must begin with an uppercase letter for json encoding.
	data := struct {
		LSB         lsb.Data         `json:"lsb"`
		Memory      memory.Data      `json:"memory"`
		CPU         cpu.Data         `json:"cpu"`
		BlockDevice blockdevice.Data `json:"block_device"`
		Languages   languages.Data   `json:"languages"`
	}{
		lsbdata,
		memdata,
		cpudata,
		blockdevicedata,
		langsdata,
	}

	b, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		fmt.Println("error:", err)
	}

	fmt.Printf("%s\n", b)
}
