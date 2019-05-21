// Copyright 2012 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"runtime"
	"time"
 
	"github.com/jpillora/present/present"
	"github.com/jpillora/present/static"
	"golang.org/x/tools/playground/socket"

	// This will register handlers at /compile and /share that will proxy to the
	// respective endpoints at play.golang.org. This allows the frontend to call
	// these endpoints without needing cross-origin request sharing (CORS).
	// Note that this is imported regardless of whether the endpoints are used or
	// not (in the case of a local socket connection, they are not called).
	_ "golang.org/x/tools/playground"
)

var scripts = []string{"jquery.js", "jquery-ui.js", "playground.js", "play.js"}

// playScript registers an HTTP handler at /play.js that serves all the
// scripts specified by the variable above, and appends a line that
// initializes the playground with the specified transport.
func playScript(router *http.ServeMux, transport string) {
	modTime := time.Now()
	var buf bytes.Buffer
	for _, p := range scripts {
		buf.Write(static.MustRead(p))
	}
	fmt.Fprintf(&buf, "\ninitPlayground(new %v());\n", transport)
	b := buf.Bytes()
	router.HandleFunc("/play.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/javascript")
		http.ServeContent(w, r, "", modTime, bytes.NewReader(b))
	})
}

func initPlayground(router *http.ServeMux, usePlay, nacl bool, originHost string) {
	if !present.PlayEnabled {
		return
	}
	if usePlay {
		playScript(router, "HTTPTransport")
		return
	}
	if nacl {
		// When specifying nativeClient, non-Go code cannot be executed
		// because the NaCl setup doesn't support doing so.
		socket.RunScripts = false
		socket.Environ = func() []string {
			if runtime.GOARCH == "amd64" {
				return environ("GOOS=nacl", "GOARCH=amd64p32")
			}
			return environ("GOOS=nacl")
		}
	}
	playScript(router, "SocketTransport")
	origin := &url.URL{Scheme: "http", Host: originHost}
	router.Handle("/socket", socket.NewHandler(origin))
}

func playable(c present.Code) bool {
	play := present.PlayEnabled && c.Play
	// TODO(jpillora): handler should be a struct, and all of these fuctions should be methods
	// Restrict playable files to only Go source files when using play.golang.org,
	// since there is no method to execute shell scripts there.
	// if *usePlayground {
	// 	return play && c.Ext == ".go"
	// }
	return play
}
