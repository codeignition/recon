// Copyright 2015 CodeIgnition. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package config

import (
	"fmt"
	"strings"
	"testing"
)

func TestGenerateUID(t *testing.T) {
	s1, err := generateUID()
	if err != nil {
		t.Error(err)
	}
	s2, err := generateUID()
	if err != nil {
		t.Error(err)
	}
	if s1 == s2 {
		t.Error("two generated UID strings are equal")
	}
}

func TestParseConfig(t *testing.T) {
	fakeUID := "23fcdd694986"
	fakeContent := fmt.Sprintf(`{"uid":"%s"}
`, fakeUID)
	r := strings.NewReader(fakeContent)
	c, err := parseConfig(r)
	if err != nil {
		t.Error(err)
	}
	if c.UID != fakeUID {
		t.Errorf("got %s; want %s", c.UID, fakeUID)
	}
}
