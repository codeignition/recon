// Copyright 2015 CodeIgnition. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package recon

// Agent represents a recon daemon running on
// a machine.
type Agent struct {
	UID string `json:"uid"`
}

type Metric struct {
	// UID of the Agent
	AgentUID string `json:"agent_uid"`

	// Metric data
	Data map[string]interface{} `json:"metrics"`
}
