// Copyright 2015 CodeIgnition. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package policy

// Type denotes the monitoring policy type.
//
// e.g. "tcp"
type Type string

// Policy is the map containing the rules of a particular monitoring policy.
//
// e.g. "tcp" PolicyType requires 2 policy keys "port" and "frequency"
type Policy struct {
	Name     string            `json:"name"`
	AgentUID string            `json:"agent_uid"`
	Type     Type              `json:"policy_type"`
	M        map[string]string `json:"m"`
}

// Config is the format used to encode/decode the monitoring policy
// received from the message queue or to store in the config file
type Config []Policy

// PolicyFuncMap maps a PolicyType to a handler function
var PolicyFuncMap = map[Type]func(Policy) error{
	"tcp": tcpPolicyHandler,
}

func (p Policy) Execute() error {
	return PolicyFuncMap[p.Type](p)
}
