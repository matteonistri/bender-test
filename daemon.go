// Package robotester provides a daemon and a simple REST API to run external
// scripts.
package main

import (
	"fmt"
	"net/http"

	"github.com/gocraft/web"
)

type Context struct {
	ScriptsDir string
}

// SetDefaults initializes Context variables
func (c *Context) SetDefaults(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	c.ScriptsDir = GetScriptsDir()
	next(w, r)
}

// RunHandler handles /run requests
func (c *Context) RunHandler(w web.ResponseWriter, r *web.Request) {
	fmt.Fprintf(w, "Requested execution of script '%s'\n", r.PathParams["script"])
}

// LogHandler handles /log requests
func (c *Context) LogHandler(w web.ResponseWriter, r *web.Request) {
	if r.PathParams["script"] != "" {
		fmt.Fprintf(w, "Requested log for script '%s'\n", r.PathParams["script"])
	} else if r.PathParams["uuid"] != "" {
		fmt.Fprintf(w, "Requested log for uuid '%s'\n", r.PathParams["uuid"])
	}
}

// StatusHandler handles /status requests
func (c *Context) StatusHandler(w web.ResponseWriter, r *web.Request) {
	if r.PathParams["script"] != "" {
		fmt.Fprintf(w, "Requested job status for script '%s\n'", r.PathParams["script"])
	} else if r.PathParams["uuid"] != "" {
		fmt.Fprintf(w, "Requested job status for uuid '%s'\n", r.PathParams["uuid"])
	} else {
		fmt.Fprintln(w, "Requested server status (general)")
		fmt.Fprintf(w, "  scripts dir: '%s'\n", c.ScriptsDir)
	}
}

const DAEMON_MODULE_NAME = "DAEMON"

func DaemonInit(address string, port string) {
	LogAppendLine(fmt.Sprintf("[%s] START", DAEMON_MODULE_NAME))

	// init http handlers
	router := web.New(Context{})
	router.Middleware((*Context).SetDefaults)
	router.Get("/run/:script", (*Context).RunHandler)
	router.Get("/log/script/:script", (*Context).LogHandler)
	router.Get("/log/uuid/:uuid", (*Context).LogHandler)
	router.Get("/status", (*Context).StatusHandler)
	router.Get("/status/script/:script", (*Context).StatusHandler)
	router.Get("/status/uuid/:uuid", (*Context).StatusHandler)

	// start http server
	LogAppendLine(fmt.Sprintf("[%s] Linsten on %s:%s", DAEMON_MODULE_NAME, address, port))
	LogFatal(http.ListenAndServe(address+":"+port, router))
}
