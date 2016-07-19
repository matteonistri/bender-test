package main

import ("time"
		"fmt")

var __SubmitChannel chan Params

type Params struct{
	name string
	uuid string
	args string
	timeout int
}

func init(){
	__SubmitChannel = make(chan Params)
	go func(){
		for {
			params := <-__SubmitChannel
			var job Job

			Run(&job, params.name, params.uuid, params.args)

			start := time.Now()
			timeout := time.Duration(params.timeout) * time.Millisecond
			for time.Since(start) <  timeout{
				fmt.Println("Reading from stdout", Log(&job))
				State(&job)
				UpdateState(job)
				if job.Status == JOB_COMPLETED{
					break
				}
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