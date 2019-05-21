// Copyright 2013 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"os"
	"net"
	"net/http"

	"github.com/jpillora/opts"
	"github.com/jpillora/present/handler"
	"github.com/jpillora/present/present"
)

type command struct {
	Host          string `opts:"help=HTTP listening interface"`
	Port          int    `opts:"help=HTTP listening port,env"`
	handler.Config
}

func main() {
	c := command{
		Host: "localhost",
		Port: 3339,
	}
	//customise defualts in google app engine
	if os.Getenv("GAE_ENV") == "standard" {
		log.Print("Configuring for App Engine Standard")
		c.Host = "0.0.0.0"
		c.UsePlayground = true
		c.ContentPath = "./content/"
	}
	//parse CLI
	opts.New(&c).Parse()
	//prepare the present handler
	h, err := handler.New(c.Config)
	if err != nil {
		log.Fatal(err)
	}
	httpAddr := fmt.Sprintf("%s:%d", c.Host, c.Port)
	ln, err := net.Listen("tcp", httpAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	if !ln.Addr().(*net.TCPAddr).IP.IsLoopback() &&
		present.PlayEnabled && !c.NativeClient && !c.UsePlayground {
		log.Print(localhostWarning)
	}
	log.Printf("Open your web browser and visit http://%s", httpAddr)
	if present.NotesEnabled {
		log.Println("Notes are enabled, press N from the browser to display them")
	}
	log.Fatal(http.Serve(ln, h))
}

const localhostWarning = `
WARNING!  WARNING!  WARNING!

The present server appears to be listening on an address that is not localhost
and is configured to run code snippets locally. Anyone with access to this address
and port will have access to this machine as the user running present.

To avoid this message, listen on localhost, run with -play=false, or run with
-play_socket=false.

If you don't understand this message, hit Control-C to terminate this process.

WARNING!  WARNING!  WARNING!
`
