package handler

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/jpillora/present/static"
	"golang.org/x/tools/present"
)

type Config struct {
	OriginHost    string `help:"host component of web origin URL"`
	ContentPath   string `help:"base path for presentation content"`
	UsePlayground bool   `help:"run code snippets using play.golang.org; if false, run them locally and deliver results by WebSocket transport"`
	NativeClient  bool   `help:"use Native Client environment playground (prevents non-Go code execution) when using local WebSocket transport"`
	PlayEnabled   bool   `help:"enable playground (permit execution of arbitrary user code)"`
	NotesEnabled  bool   `help:"enable presenter notes (press 'N' from the browser to display them)"`
}

func New(c Config) (http.Handler, error) {
	// terrible globals
	if c.PlayEnabled {
		present.PlayEnabled = true
	}
	if c.NotesEnabled {
		present.NotesEnabled = true
	}
	if c.OriginHost == "" {
		c.OriginHost = "localhost"
	}
	if c.ContentPath == "" {
		pwd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("Couldn't get pwd: %s", err)
		}
		c.ContentPath = pwd
	}
	err := initTemplates()
	if err != nil {
		return nil, fmt.Errorf("Failed to parse templates: %v", err)
	}
	router := http.NewServeMux()
	initPlayground(router, c.UsePlayground, c.NativeClient, c.OriginHost)
	router.Handle("/static/", http.StripPrefix("/static/", static.Handler()))
	router.Handle("/", dirHandler(c.ContentPath))
	return router, nil
}

func environ(vars ...string) []string {
	env := os.Environ()
	for _, r := range vars {
		k := strings.SplitAfter(r, "=")[0]
		var found bool
		for i, v := range env {
			if strings.HasPrefix(v, k) {
				env[i] = r
				found = true
			}
		}
		if !found {
			env = append(env, r)
		}
	}
	return env
}
