// Copyright 2015 CodeIgnition. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package handlers

import "github.com/codeignition/recon/policy"

func init() {
	policy.RegisterHandler("tcp", TCP)
	policy.RegisterHandler("system_data", SystemData)
}
