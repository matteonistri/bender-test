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
	ip      string
}

//Receive a job from channel and call the runner to execute it
func workerLoop() {
	for {
		params := <-submitChannel

		job := &Job{}
		ret := job.Run(params.name, params.uuid, params.ip, params.args)
		if ret < 0 {
			job.Status = JobNotFound
		} else {
			logChannel := job.Log()
			stateChannel := job.State()
			previousState := ""

			var exit bool
			rep := &ReportContext{}
			err := rep.New(params.name, params.uuid, time.Now(), true)
			if err != nil {
				LogErr(logContextWorker, "Error while creating report: %s", err.Error())
				exit = true
			} else {
				exit = false
			}

			for !exit {
				select {
				case m := <-logChannel:
					LogDeb(logContextWorker, m)
					rep.UpdateString(m)
					wd := WebData{
						Msg: m,
						Datatype: "output",
						Ip: job.ip,
					}
					webChannel <- wd
					time.Sleep(500 * time.Millisecond)
				case s := <-stateChannel:
					if previousState != s {
						LogDeb(logContextWorker, "Receive [%v] state [%v]", job.Name, s)
						job.Status = s
						previousState = s
						wd := WebData{
							Msg: s,
							Datatype: "scriptstatus",
							Ip: job.ip,
						}
						webChannel <- wd
					}
					if s != JobWorking {
						LogInf(logContextWorker, "%v", job)
						exit = true
					}
					workerLocalStatus.SetState(*job)
				case <-time.After(params.timeout * time.Second):
					LogDeb(logContextWorker, "Exec script [%v] Timeout! [%v]", job.Name, params.timeout*time.Second)
					LogInf(logContextWorker, "%v", job)
					job.Status = JobTimeout
					wd := WebData{
						Msg: JobTimeout,
						Datatype: "scriptstatus",
						Ip: job.ip,
					}
					webChannel <- wd
					exit = true
					workerLocalStatus.SetState(*job)
				}
			}
		}
	}
}

//Submit send a new job on the channel
func Submit(name, uuid, ip string, argsMap url.Values, timeout time.Duration) {
	var args []string
	for k, v := range argsMap {
		for _, x := range v {
			args = append(args, k)
			if x != "" {
				args = append(args, string(x))
			}
		}
	}

	params := params{
		name:    name,
		uuid:    uuid,
		args:    args,
		timeout: timeout,
		ip:      ip,
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
