package main

import ("time"
		"fmt")

var __SubmitChannel chan Params
var logContextWorker LoggerContext

type Params struct{
	name string
	uuid string
	args string
	timeout int
}

func init(){
	logContextWorker = LoggerContext{
		name: "WORKER",
		level: 3}

	__SubmitChannel = make(chan Params)
	go func(){
		for {
			params := <-__SubmitChannel
			var job Job

			ret := Run(&job, params.name, params.uuid, params.args)
			logChan := *Log()

			if ret == 0 {
				start := time.Now()
				timeout := time.Duration(params.timeout) * time.Millisecond
				for time.Since(start) <  timeout{
					State(&job)
					UpdateState(job)
					if job.Status != JOB_WORKING{
						break
					}
					select{
						case out := <-logChan:
							fmt.Println(out)
						default:
							time.Sleep(500 * time.Millisecond)
					}
				}

				if time.Since(start) >  timeout{
	 				LogWar(logContextWorker, "Execution timed out")
	 				job.Status = JOB_FAILED
				}
			} else {
				job.Status = JOB_NOT_FOUND
			}
        }
    }()
}

func Submit(name, uuid, args string, timeout int){
	params := Params{
		name:    name,
		uuid:    uuid,
		args:    args,
		timeout: timeout}

	__SubmitChannel <- params
}

func UpdateState(job Job){
	fmt.Println(job.Status)
}