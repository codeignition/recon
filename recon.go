// Copyright 2015 Hari haran. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"log"

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

	data := struct {
		Lsb    lsb.Data    `json:"lsb"`
		Memory memory.Data `json:"memory"`
	}{
		lsbdata,
		memdata,
	}
	b, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("%s\n", b)
}
