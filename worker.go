package main

import (
	"fmt"
	"net/url"
	"time"
)

var __SubmitChannel chan Params
var logContextWorker LoggerContext

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
				for time.Since(start) < timeout {
					select {
					case out := <-logChan:
						fmt.Println(out)
					default:
						time.Sleep(20 * time.Millisecond)
					}
					State(&job)
					//UpdateState(job)
					if job.Status != JOB_WORKING {
						break
					}
				}

				if time.Since(start) > timeout {
					LogWar(logContextWorker, "Execution timed out")
					job.Status = JOB_FAILED
				}
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

func UpdateState(job Job) {
	fmt.Println(job.Status)
}
