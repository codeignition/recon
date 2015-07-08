// Copyright 2015 Hari haran. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func reconHandler(w http.ResponseWriter, r *http.Request) {
	// Adding the charset parameter explicitly makes browsers show the unicode characters properly.
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	data := accumulateData()

	var (
		b   []byte
		err error
	)
	if r.FormValue("indent") == "1" {
		b, err = json.MarshalIndent(data, "", "    ")
	} else {
		b, err = json.Marshal(data)
	}

	if err != nil {
		log.Println(err)

		w.WriteHeader(http.StatusInternalServerError)

		// We don't want to show the actual error to the user, for
		// security and user experience. So we show a custom error message.
		fmt.Fprintf(w, `{"error":"%s"}`, "unable to marshal JSON")

		// Terminate with a newline to make the output look a little
		// nicer when debugging.
		w.Write([]byte("\n"))
		return
	}

	fmt.Fprintf(w, "%s\n", b)
}
