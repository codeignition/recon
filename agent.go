// Copyright 2015 CodeIgnition. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package recon

// Agent represents a recon daemon running on
// a machine.
type Agent struct {
	UID      string `json:"uid"`
	HostName string `json:"host_name"`
}
