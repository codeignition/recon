// Copyright 2015 Hari haran. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"os/user"
	"time"

	"github.com/hariharan-uno/recon/blockdevice"
	"github.com/hariharan-uno/recon/cpu"
	"github.com/hariharan-uno/recon/etc"
	"github.com/hariharan-uno/recon/initpackage"
	"github.com/hariharan-uno/recon/kernel"
	"github.com/hariharan-uno/recon/languages"
	"github.com/hariharan-uno/recon/lsb"
	"github.com/hariharan-uno/recon/memory"
	"github.com/hariharan-uno/recon/network"
	"github.com/hariharan-uno/recon/ps"
	"github.com/hariharan-uno/recon/uptime"
)

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
	kerneldata, err := kernel.CollectData()
	if err != nil {
		log.Println(err)
	}
	currentUser, err := user.Current()
	if err != nil {
		log.Println(err)
	}
	etcdata, err := etc.CollectData()
	if err != nil {
		log.Println(err)
	}
	netdata, err := network.CollectData()
	if err != nil {
		log.Println(err)
	}
	data := map[string]interface{}{
		"lsb":          lsbdata,
		"memory":       memdata,
		"cpu":          cpudata,
		"block_device": blockdevicedata,
		"languages":    langsdata,
		"kernel":       kerneldata,
		"recon_time":   time.Now(),
		"init_package": initpackage.Name,
		"command":      map[string]string{"ps": ps.Command},
		"current_user": currentUser.Username, // if more data is required, use currentUser instead of just the Username field
		"etc":          etcdata,
		"network":      netdata,
	}
	copyMap(uptimedata, data) // uptime Data is not namespaced.
	return data
}
