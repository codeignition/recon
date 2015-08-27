// Copyright 2015 CodeIgnition. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package policy

import (
	"errors"

	"golang.org/x/net/context"
)

// Type denotes the monitoring policy type.
//
// e.g. "tcp"
type Type string

// Policy is the map containing the rules of a particular monitoring policy.
//
// e.g. "tcp" PolicyType requires 2 policy keys "port" and "frequency"
type Policy struct {
	Name     string
	AgentUID string
	Type     Type
	M        map[string]string
}

// Config is the format used to encode/decode the monitoring policy
// received from the message queue or to store in the config file
type Config []Policy

// policyFuncMap maps a PolicyType to a handler function
var policyFuncMap = map[Type]func(context.Context, Policy) (<-chan Event, error){
	"tcp": tcpPolicyHandler,
}

func (p Policy) Execute() (<-chan Event, error) {
	if err := p.Valid(); err != nil {
		return nil, err
	}
	return policyFuncMap[p.Type](context.TODO(), p)
}

// Valid checks whether the policy is valid.
func (p Policy) Valid() error {
	if p.Name == "" {
		return errors.New("policy name can't be empty")
	}
	if _, ok := policyFuncMap[p.Type]; !ok {
		return errors.New("policy type unknown")
	}
	return nil
}
