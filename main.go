package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/jpillora/cookieauth"
	"github.com/jpillora/opts"
	"github.com/jpillora/present/handler"
	"github.com/jpillora/requestlog"
	"golang.org/x/tools/present"
)

type command struct {
	Host string `opts:"help=HTTP listening interface"`
	Port int    `opts:"help=HTTP listening port, env"`
	Log  bool   `opts:"help=enable request logging"`
	Auth string `opts:"help=enable basic-auth (user:pass)"`
	handler.Config
}

func main() {
	c := command{
		Host: "localhost",
		Port: 3999,
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
	//middleware
	if a := strings.Split(c.Auth, ":"); len(a) == 2 {
		h = cookieauth.Wrap(h, a[0], a[1])
	}
	if c.Log {
		h = requestlog.Wrap(h)
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
