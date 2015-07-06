// Copyright 2015 Hari haran. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"

	"github.com/hariharan-uno/recon/lsb"
)

func main() {
	log.SetPrefix("recon: ")
	d, err := lsb.CollectData()
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("%v\n", d)
}
