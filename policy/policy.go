// Copyright 2015 CodeIgnition. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package policy

import (
	"errors"
	"sync"

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
var handlerFuncMap = struct {
	sync.Mutex
	m map[string]HandlerFunc
}{
	m: make(map[string]HandlerFunc),
}

func (p Policy) Execute(ctx context.Context) (<-chan Event, error) {
	if err := p.Valid(); err != nil {
		return nil, err
	}

	handlerFuncMap.Lock()
	f := handlerFuncMap.m[p.Type]
	handlerFuncMap.Unlock()

	return f(ctx, p)
}

// Valid checks whether the policy is valid.
func (p Policy) Valid() error {
	if p.Name == "" {
		return errors.New("policy name can't be empty")
	}

	handlerFuncMap.Lock()
	_, ok := handlerFuncMap.m[p.Type]
	handlerFuncMap.Unlock()

	if !ok {
		return errors.New("policy type unknown")
	}
	return nil
}

func RegisterHandler(policyType string, handlerFunc HandlerFunc) error {
	if policyType == "" {
		return errors.New("policy type can't be empty")
	}

	handlerFuncMap.Lock()
	defer handlerFuncMap.Unlock()

	if _, ok := handlerFuncMap.m[policyType]; ok {
		return errors.New("handler for the policy type already exists")
	}

	handlerFuncMap.m[policyType] = handlerFunc
	return nil
}
