// Copyright 2015 CodeIgnition. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package policy

import (
	"errors"

	"golang.org/x/net/context"
)

// Policy represents a monitoring policy
type Policy struct {
	Name     string            // Name of the monitoring policy
	AgentUID string            // Agent UID
	Type     string            // Type denotes the monitoring policy type. e.g. "tcp"
	M        map[string]string // M is the map containing the rules of a particular monitoring policy.
}

// Config is a slice of monitoring policies
type Config []Policy

// HandlerFunc is the type of a policy handler function. Any policy handler
// function must be of this type.
type HandlerFunc func(context.Context, Policy) (<-chan Event, error)

// handlerFuncMap maps a policy type to a handler function
var handlerFuncMap = make(map[string]HandlerFunc)

func (p Policy) Execute(ctx context.Context) (<-chan Event, error) {
	if err := p.Valid(); err != nil {
		return nil, err
	}
	return handlerFuncMap[p.Type](ctx, p)
}

// Valid checks whether the policy is valid.
func (p Policy) Valid() error {
	if p.Name == "" {
		return errors.New("policy name can't be empty")
	}
	if _, ok := handlerFuncMap[p.Type]; !ok {
		return errors.New("policy type unknown")
	}
	return nil
}

func RegisterHandler(policyType string, handlerFunc HandlerFunc) error {
	if policyType == "" {
		return errors.New("policy type can't be empty")
	}

	if _, ok := handlerFuncMap[policyType]; ok {
		return errors.New("handler for the policy type already exists")
	}
	handlerFuncMap[policyType] = handlerFunc
	return nil
}
