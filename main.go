// Copyright 2015 Hari haran. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	log.SetPrefix("recon: ")

	var addr = flag.String("addr", ":3030", "serve HTTP on `address`")
	flag.Parse()

	http.HandleFunc("/", reconHandler)
	fmt.Printf("recon: starting the server on http://localhost%s\n", *addr)
	fmt.Printf("recon: if you'd like the JSON to be indented, append %q to the above URL\n", "?indent=1")
	log.Fatal(http.ListenAndServe(*addr, nil))
}
