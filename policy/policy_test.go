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
	interval, ok := p.M["interval"]
	if !ok {
		return nil, errors.New(`"interval" key missing in fake policy`)
	}
	d, err := time.ParseDuration(interval)
	if err != nil {
		return nil, err
	}

	// This check is here to ensure time.Ticker(d) doesn't panic
	if d <= 0 {
		return nil, errors.New("frequency must be a positive quantity")
	}

	out := make(chan Event)
	go func() {
		t := time.NewTicker(d)
		for {
			select {
			case <-ctx.Done():
				t.Stop()
				close(out)
				return
			case <-t.C:
				out <- Event{
					Time:   time.Now(),
					Policy: p,
					Data: map[string]interface{}{
						"foo": foo,
					},
				}
			}
		}
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

func TestExecute(t *testing.T) {
	p := Policy{
		Name: "dummy",
		Type: "fake",
		M: map[string]string{
			"foo":      "foo_value",
			"interval": "200ms",
		},
	}
	ctx, cancel := context.WithCancel(context.Background())
	out, err := p.Execute(ctx)
	if err != nil {
		t.Error(err)
	}
	go func() {
		time.Sleep(1 * time.Second)
		cancel()
	}()

	var count int
	// Here, we are also able to test whether out is being
	// closed when cancel() is called. If out is not closed,
	// this test should have been running forever.
	for evt := range out {
		count++
		if evt.Data["foo"] != "foo_value" {
			t.Errorf(`want evt.Data["foo"] = %s; got %s`, "foo", evt.Data["foo"])
		}
	}

	// The interval for the dummy policy is 200ms.
	// We are calling cancel after 1 sec. Typically, we receive
	// 4 or 5 events in that duration.
	if count != 4 && count != 5 {
		t.Errorf(`want count to be either 4 or 5; got %d`, count)
	}
}
