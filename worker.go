package main

import (
	"net/url"
	"time"
)

var submitChannel chan params
var logContextWorker LoggerContext
var workerLocalStatus *StatusModule
var endReadLog = make(chan bool)

type params struct {
	name    string
	uuid    string
	args    []string
	timeout int
}

//Receive a job from channel and call the runner to execute it
func init() {
	logContextWorker = LoggerContext{
		name:  "WORKER",
		level: 3}

	submitChannel = make(chan params)
	go func() {
		for {
			params := <-submitChannel
			job := &Job{}
			ret := job.Run(params.name, params.uuid, params.args)

			if ret == 0 {
				start := time.Now()
				timeout := time.Duration(params.timeout) * time.Millisecond

				rep := &ReportContext{}
				rep.New(params.name, params.uuid, start, true)
				logChan := *Log()

			timeToLive:
				for time.Since(start) < timeout {
					select {
					case m := <-logChan:
						rep.UpdateString(m)
					case <-endReadLog:
						LogDeb(logContextWorker, "received end of read sync")
						break timeToLive
					default:
						time.Sleep(20 * time.Millisecond)
					}
					job.State()
					workerLocalStatus.SetState(*job)
				}

				if time.Since(start) > timeout {
					LogWar(logContextWorker, "Execution timed out")
					job.Status = JobFailed
				}

				workerLocalStatus.SetState(*job)
			} else {
				job.Status = JobNotFound
			}
		}
	}()
}

//Submit send a new job on the channel
func Submit(name, uuid string, argsMap url.Values, timeout int) {
	var args []string
	for k, v := range argsMap {
		for _, x := range v {
			args = append(args, k)
			args = append(args, string(x))
		}
	}

	params := params{
		name:    name,
		uuid:    uuid,
		args:    args,
		timeout: timeout}

	submitChannel <- params
}

// WorkerInit ...
func WorkerInit(sm *StatusModule) {
	workerLocalStatus = sm
}
