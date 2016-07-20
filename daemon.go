// Package robotester provides a daemon and a simple REST API to run external
// scripts.
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gocraft/web"
	"github.com/satori/go.uuid"
)

type statusDaemon struct {
	ServerStatus serverStatus `json:"serverStatus"`
	ServerName   string       `json:"serverName"`
	Timestamp    time.Time    `json:"timestamp"`
}

type statusJobs struct {
	Jobs []Job `json:"jobs"`
}

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
	r.ParseForm()

	name := r.PathParams["script"]
	uuid := uuid.NewV4()
	timeout := 600
	args := r.Form

	if sd.ServerStatus == SERVER_WORKING {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
	Submit(name, uuid, args, timeout)
}

// LogHandler handles /log requests
func (c *Context) LogHandler(w web.ResponseWriter, r *web.Request) {
	if r.PathParams["script"] != "" {
		fmt.Fprintf(w, "Requested log for script '%s'\n", r.PathParams["script"])
	} else if r.PathParams["uuid"] != "" {
		fmt.Fprintf(w, "Requested log for uuid '%s'\n", r.PathParams["uuid"])
	}
}

// StatusHandler handles /state requests
func (c *Context) StatusHandler(w web.ResponseWriter, r *web.Request) {
	//general state requests
	if r.RequestURI == "/state" {
		js, err := json.Marshal(sd)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			panic("json creation failed")
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	} else {
		// script-name specific requests
		r.ParseForm()
		response := statusJobs{
			Jobs: sm.GetJobs(r.PathParams["script"])}
		js, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			panic("json creation failed")
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

// InitDaemonStatus initializes a stausDaemon struct
func InitDaemonStatus(serverName string) statusDaemon {
	sd = statusDaemon{
		ServerStatus: SERVER_IDLE,
		ServerName:   serverName,
		Timestamp:    time.Now()}

	return sd
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
	router.Get("/state", (*Context).StatusHandler)
	router.Get("/state/:script", (*Context).StatusHandler)

	// start http server
	LogAppendLine(fmt.Sprintf("[%s] Linsten on %s:%s", DAEMON_MODULE_NAME, address, port))
	LogFatal(http.ListenAndServe(address+":"+port, router))
}
