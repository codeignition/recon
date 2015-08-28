// Copyright 2015 CodeIgnition. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package policy

import (
	"errors"
	"testing"
	"time"

	"golang.org/x/net/context"
)

func fakePolicyHandler(ctx context.Context, p Policy) (<-chan Event, error) {
	foo, ok := p.M["foo"]
	if !ok {
		return nil, errors.New(`"foo" key missing in fake policy`)
	}
	bar, ok := p.M["bar"]
	if !ok {
		return nil, errors.New(`"bar" key missing in fake policy`)
	}
	out := make(chan Event)
	go func() {
		out <- Event{
			Time:   time.Now(),
			Policy: p,
			Data: map[string]interface{}{
				"count": 1,
				"foo":   foo,
				"bar":   bar,
			},
		}
		out <- Event{
			Time:   time.Now(),
			Policy: p,
			Data: map[string]interface{}{
				"count": 2,
				"foo":   foo,
				"bar":   bar,
			},
		}
		close(out)
	}()
	return out, nil
}

func TestNewHandler(t *testing.T) {
	err := NewHandler("", fakePolicyHandler)
	if err == nil {
		t.Fatal(errors.New("NewHandler should return an error when the type is empty"))
	}
	err = NewHandler("fake", fakePolicyHandler)
	if err != nil {
		t.Fatal(err)
	}
}

func TestValid(t *testing.T) {
	p := Policy{
		Name: "",
	}
	if err := p.Valid(); err == nil {
		t.Error(`want error "policy name can't be empty" got nil`)
	}
	p = Policy{
		Name: "dummy",
		Type: "unknownDummyType",
	}
	if err := p.Valid(); err == nil {
		t.Error(`want error "policy type unknown" got nil`)
	}
}
