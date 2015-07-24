// Copyright 2015 CodeIgnition. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

// Package fileutil provides utility functions for dealing with files.
package fileutil

import (
	"os"
)

// Exists returns true if a file exists.
func Exists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}
