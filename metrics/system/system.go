// Copyright 2015 CodeIgnition. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

// +build linux

// Package system provides system metrics.
package system

import (
	"github.com/codeignition/recon/metrics/system/disk"
	"github.com/codeignition/recon/metrics/system/top"
)

// Data denotes system data
type Data map[string]interface{}

// Merge merges the input data with caller.
// It will overwrite the value of a key if it exists already.
func (d Data) Merge(a map[string]interface{}) {
	for k, v := range a {
		d[k] = v
	}
}

// CollectData collects the data and returns an error if any.
func CollectData() (Data, error) {
	d := make(Data)
	top, err := top.CollectData()
	if err != nil {
		return d, err
	}
	d.Merge(top)
	disk, err := disk.CollectData()
	if err != nil {
		return d, err
	}
	d.Merge(disk)
	return d, nil
}
