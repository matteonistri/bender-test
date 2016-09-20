// Package robotester provides a daemon and a simple REST API to run external
// scripts.
package main

import (
	"encoding/json"
	"net/http"
	"time"
	"os"
	"io/ioutil"

	"github.com/gocraft/web"
	"github.com/satori/go.uuid"
	"github.com/gorilla/websocket"
)

var logContextDaemon LoggerContext
var daemonLocalStatus *StatusModule

type statusJobs struct {
	Jobs []Job `json:"jobs"`
}

// Context ...
type Context struct {
	ScriptsDir string
}

// SetDefaults initializes Context variables
func (c *Context) SetDefaults(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	next(w, r)
}

// RunHandler handles /run requests
func (c *Context) RunHandler(w web.ResponseWriter, r *web.Request) {
	LogInf(logContextDaemon, "Receive RUN[%v] request from: %v", "Daemon", r.RemoteAddr)
	r.ParseForm()

	name := r.PathParams["script"]
	uuid := uuid.NewV4().String()
	timeout := 10
	params := r.Form

	status, _ := daemonLocalStatus.GetState()
	if status == DaemonWorking {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
	js, err := json.Marshal(uuid)

	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		LogErr(logContextDaemon, "json creation failed")
		return
	}

	w.Write(js)

	Submit(name, uuid, params, time.Duration(timeout))
}

// LogHandler handles /log requests
func (c *Context) LogHandler(w web.ResponseWriter, r *web.Request) {
	//LogInf(logContextDaemon, "Receive LOG[%v] request from: %v", "Daemon", r.RemoteAddr)
	r.ParseForm()
	name := r.PathParams["script"]
	ids := r.Form["uuid"]
	var buffer []byte
	var list [][]string
	var js []byte
	var err error

	if len(ids) > 0 {
		rep := GetReportContext()
		buffer = rep.Read(rep.offset, -1)
		rep.offset += int64(len(buffer))

		js, err = json.Marshal(string(buffer))
	} else {
		list, err = ReportList(name)
		if err != nil {
			LogErr(logContextDaemon, "Unable to find log for this script")
			return
		}
		js, err = json.Marshal(list)
	}

	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		LogErr(logContextDaemon, "json creation failed")
		return
	}

	w.Write(js)
}

// StatusHandler handles /state requests
func (c *Context) StatusHandler(w web.ResponseWriter, r *web.Request) {
	//general state requests

	if r.RequestURI == "/state" {
		//LogInf(logContextDaemon, "Receive STATE[%v] request from: %v", "Daemon", r.RemoteAddr)
		js, err := json.Marshal(daemonLocalStatus)

		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			LogErr(logContextDaemon, "json creation failed")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	} else {
		// script-name specific requests
		r.ParseForm()

		LogInf(logContextDaemon, "Receive STATE[%v] request from: %v", r.PathParams["script"], r.RemoteAddr)

		response := statusJobs{
			Jobs: daemonLocalStatus.GetJobs(r.PathParams["script"])}
		js, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			LogErr(logContextDaemon, "json creation failed")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

// HomeHandler ...
func (c *Context) HomeHandler(w web.ResponseWriter, r *web.Request) {
	LogInf(logContextDaemon, "Receive HOME[%v] request from: %v", "Daemon", r.RemoteAddr)

	indexFile, err := os.Open("static/run.html")
    if err != nil {
        LogErr(logContextDaemon, "cannot open html home file")
        return
    }
    index, err := ioutil.ReadAll(indexFile)
    if err != nil {
        LogErr(logContextDaemon, "cannot read from html home file")
        return
    }
    w.Write(index)
}

//ListHandler ...
func (c *Context) ListHandler(w web.ResponseWriter, r *web.Request) {
	LogInf(logContextDaemon, "Receive LIST[%v] request from: %v", "Daemon", r.RemoteAddr)
	scripts := List()
	js, err := json.Marshal(scripts)

	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		LogErr(logContextDaemon, "json creation failed")
		return
	}

	w.Write(js)
}

// SetListHandler ...
func (c *Context) SetListHandler(w web.ResponseWriter, r *web.Request) {
	LogInf(logContextDaemon, "Receive SETS[%v] request from: %v", "Daemon", r.RemoteAddr)
	r.ParseForm()
	list := r.Form["set"]
	var js []byte
	var err error

	if len(list) <= 0 {
		sets := SetsList()
		js, err = json.Marshal(sets)
	} else {
		set := GetSet(list[0])
		js, err = json.Marshal(set)
	}

	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		LogErr(logContextDaemon, "json creation failed")
		return
	}

	w.Write(js)
}

func (c *Context) Websocket(w web.ResponseWriter, r *web.Request) {
	LogInf(logContextDaemon, "Receive WEBSOCKET[%v] request from: %v", "Daemon", r.RemoteAddr)
	var upgrader = websocket.Upgrader {
    	ReadBufferSize: 1024,
    	WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(w, r.Request, nil)
  	if err != nil {
    	LogErr(logContextDaemon, "websocket connection establishment failed")
    	return
  	}
	LogInf(logContextDaemon,"Websocket connected")

	previous := ""
  	for {
  		current, _ := daemonLocalStatus.GetState()

  		if current != previous {
  			js, err :=  json.Marshal(current)
	  		if err != nil {
		        LogErr(logContextDaemon, "json creation failed")
		    }

		    err = conn.WriteMessage(websocket.TextMessage, js)
		    if err != nil {
		        LogErr(logContextDaemon, "websocket message sending failed")
		    }
		    previous = current
  		}
  	}
}

// DaemonInit ...
func DaemonInit(sm *StatusModule, cm *ConfigModule) {
	daemonLocalStatus = sm

	// init logger
	logContextDaemon = LoggerContext{
		level: cm.GetLogLevel("daemon", 3),
		name:  "DAEMON"}
	LogInf(logContextDaemon, "START")

	// init http handlers
	router := web.New(Context{})
	router.Middleware((*Context).SetDefaults)
	router.Middleware(web.StaticMiddleware("static"))
	router.Get("/run/:script", (*Context).RunHandler)
	router.Get("/log/:script", (*Context).LogHandler)
	router.Get("/state", (*Context).StatusHandler)
	router.Get("/state/:script", (*Context).StatusHandler)
	router.Get("/", (*Context).HomeHandler)
	router.Get("/service/list", (*Context).ListHandler)
	router.Get("/service/sets", (*Context).SetListHandler)
	router.Get("/websocket", (*Context).Websocket)

	// start http server
	address := cm.Get("daemon", "address", "0.0.0.0")
	port := cm.Get("daemon", "port", "8080")
	LogInf(logContextDaemon, "Listening on %s:%s", address, port)
	LogFatal(logContextDaemon, http.ListenAndServe(address+":"+port, router))
}
