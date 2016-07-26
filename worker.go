package main

import (
	"fmt"
	"net/url"
	"time"
)

var __SubmitChannel chan Params
var logContextWorker LoggerContext
var worker_localStatus *StatusModule
var endReadLog = make(chan bool)

type Params struct {
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

	__SubmitChannel = make(chan Params)
	go func() {
		for {
			params := <-__SubmitChannel
			var job Job

			ret := Run(&job, params.name, params.uuid, params.args)
			logChan := *Log()

			if ret == 0 {
				start := time.Now()
				timeout := time.Duration(params.timeout) * time.Millisecond
			timeToLive:
				for time.Since(start) < timeout {
					select {
					case out := <-logChan:
						fmt.Println(out)
					case <-endReadLog:
						LogDeb(logContextWorker, "received end of read sync")
						break timeToLive
					default:
						time.Sleep(20 * time.Millisecond)
					}
					State(&job)
					worker_localStatus.SetState(job)
				}

				if time.Since(start) > timeout {
					LogWar(logContextWorker, "Execution timed out")
					job.Status = JOB_FAILED
				}
				worker_localStatus.SetState(job)
			} else {
				job.Status = JOB_NOT_FOUND
			}
		}
	}()
}

//Send a new job on the channel
func Submit(name, uuid string, argsMap url.Values, timeout int) {
	var args []string
	for k, v := range argsMap {
		for _, x := range v {
			args = append(args, k)
			args = append(args, string(x))
		}
	}

	params := Params{
		name:    name,
		uuid:    uuid,
		args:    args,
		timeout: timeout}

	__SubmitChannel <- params
}

func WorkerInit(sm *StatusModule) {
	worker_localStatus = sm
}
