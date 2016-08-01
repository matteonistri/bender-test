package main

import (
	"net/url"
	"time"
)

var submitChannel chan params
var logContextWorker LoggerContext
var workerLocalStatus *StatusModule

type params struct {
	name    string
	uuid    string
	args    []string
	timeout time.Duration
}

//Receive a job from channel and call the runner to execute it
func workerLoop() {
	for {
		params := <-submitChannel

		job := &Job{}
		ret := job.Run(params.name, params.uuid, params.args)
		if ret < 0 {
			job.Status = JobNotFound
		} else {
			logChannel := job.Log()
			stateChannel := job.State()
			previousState := ""

			//rep := &ReportContext{}
			//rep.New(params.name, params.uuid, start, true)
			exit := false
			for !exit {
				select {
				case m := <-logChannel:
					LogDeb(logContextWorker, m)
				case s := <-stateChannel:
					if previousState != s {
						LogDeb(logContextWorker, "Receive [%v] state [%v]", job.Name, s)
						previousState = s
					}
					if s != JobWorking {
						LogInf(logContextWorker, "%v", job)
						exit = true
					}
				case <-time.After(params.timeout * time.Second):
					LogDeb(logContextWorker, "Exec script [%v] Timeout! [%v]", job.Name, params.timeout*time.Second)
					LogInf(logContextWorker, "%v", job)
					job.Name = JobTimeout
					exit = true
				}
			}
		}
		workerLocalStatus.SetState(*job)
	}
}

//Submit send a new job on the channel
func Submit(name, uuid string, argsMap url.Values, timeout time.Duration) {
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
		timeout: timeout,
	}

	submitChannel <- params
}

// WorkerInit ...
func WorkerInit(sm *StatusModule) {
	workerLocalStatus = sm
	logContextWorker = LoggerContext{
		name:  "WORKER",
		level: 3}

	submitChannel = make(chan params)

	go workerLoop()
}
