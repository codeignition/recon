// Copyright 2015 CodeIgnition. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package policy

import (
	"errors"

	"golang.org/x/net/context"
)

// Policy is the map containing the rules of a particular monitoring policy.
//
// e.g. "tcp" PolicyType requires 2 policy keys "port" and "frequency"
type Policy struct {
	Name     string
	AgentUID string
	Type     string // Type denotes the monitoring policy type. e.g. "tcp"
	M        map[string]string
}

// Config is the format used to encode/decode the monitoring policy
// received from the message queue or to store in the config file
type Config []Policy

type HandlerFunc func(context.Context, Policy) (<-chan Event, error)

// policyFuncMap maps a policy type to a handler function
var policyFuncMap = make(map[string]HandlerFunc)

func (p Policy) Execute(ctx context.Context) (<-chan Event, error) {
	if err := p.Valid(); err != nil {
		return nil, err
	}
	return policyFuncMap[p.Type](ctx, p)
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

func RegisterHandler(policyType string, handlerFunc HandlerFunc) error {
	if policyType == "" {
		return errors.New("policy type can't be empty")
	}

	if _, ok := policyFuncMap[policyType]; ok {
		return errors.New("handler for the policy type already exists")
	}
	policyFuncMap[policyType] = handlerFunc
	return nil
}
