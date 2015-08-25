// Copyright 2015 CodeIgnition. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/codeignition/recon/internal/fileutil"
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
	fakeContent := fmt.Sprintf(`{"UID":"%s"}
`, fakeUID)
	r := strings.NewReader(fakeContent)
	c, err := parseConfig(r)
	if err != nil {
		t.Error(err)
	}
	if c.UID != fakeUID {
		t.Errorf("got config UID %q; want %q", c.UID, fakeUID)
	}
}

func TestInitExisting(t *testing.T) {
	f, err := ioutil.TempFile("", "recond_fake_config")
	if err != nil {
		t.Error(err)
	}
	fakeUID := "13fcdf794886"
	fakeContent := fmt.Sprintf(`{"UID":"%s"}
`, fakeUID)
	_, err = f.WriteString(fakeContent)
	if err != nil {
		t.Error(err)
	}
	if err := f.Close(); err != nil {
		t.Error(err)
	}
	configPath = f.Name()
	c, err := Init()
	if err != nil {
		t.Error(err)
	}
	if c.UID != fakeUID {
		t.Errorf("got config UID %q; want %q", c.UID, fakeUID)
	}
	if err := os.Remove(f.Name()); err != nil {
		t.Error(err)
	}
}

func TestInitNew(t *testing.T) {
	configPath = filepath.Join(os.TempDir(), "recond_fake_new_config") // doesn't exist
	c, err := Init()
	if err != nil {
		t.Error(err)
	}
	if !fileutil.Exists(configPath) {
		t.Errorf("Init didn't create the config file")
	}
	if c.UID == "" {
		t.Errorf("got config UID as an empty string; want a non empty string")
	}
	if err := os.Remove(configPath); err != nil {
		t.Error(err)
	}
}

func TestSave(t *testing.T) {
	f, err := ioutil.TempFile("", "recond_fake_config")
	if err != nil {
		t.Error(err)
	}
	fakeUID := "13fcdf794886"
	fakeContent := fmt.Sprintf(`{"UID":"%s"}
`, fakeUID)
	_, err = f.WriteString(fakeContent)
	if err != nil {
		t.Error(err)
	}
	if err := f.Close(); err != nil {
		t.Error(err)
	}
	configPath = f.Name()
	c, err := Init()
	if err != nil {
		t.Error(err)
	}
	// this is redundant as it is checked in TestInit but lets leave it anyways
	if c.UID != fakeUID {
		t.Errorf("got config UID %q; want %q", c.UID, fakeUID)
	}
	s, err := generateUID()
	if err != nil {
		t.Error(err)
	}
	c.UID = s
	if err := c.Save(); err != nil {
		t.Error(err)
	}
	c, err = Init()
	if err != nil {
		t.Error(err)
	}
	if c.UID != s {
		t.Errorf("got config UID %q; want %q", c.UID, s)
	}
	if err := os.Remove(f.Name()); err != nil {
		t.Error(err)
	}
}
