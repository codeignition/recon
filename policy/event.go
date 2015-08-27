// Copyright 2015 CodeIgnition. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package policy

import "time"

// Event data that will be sent by policy handlers
type Event struct {
	Time   time.Time
	Policy Policy
	Data   map[string]interface{} // Data may include status, stats, etc.
}
