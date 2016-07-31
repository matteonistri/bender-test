package main

import (
	"fmt"
	"testing"
	"time"
)

func runner(name string) {
	job := &Job{}
	ret := job.Run(name, "f3628fa3-2b25-4dfd-ac2d-2c8a8613915c", []string{})
	fmt.Println(ret)
	logChannel := *job.Log()
	stateChannel := *job.State()

	exit := true
	previousState := ""
	for exit {
		select {
		case m := <-logChannel:
			//rep.UpdateString(m)
			fmt.Println(m)
		case s := <-stateChannel:
			if previousState != s {
				LogDeb(logContextTestRunner, "received state [%v]", s)
				previousState = s
			}
			if s == JobCompleted || s == JobFailed {
				exit = false
			}
			break
		case <-time.After(60 * time.Second):
			LogDeb(logContextTestRunner, "Timeout!")
			exit = false
			break
		}
	}
}

var logContextTestRunner LoggerContext

func TestRunnerF(t *testing.T) {
	var cm ConfigModule
	ConfigInit(&cm, "bender-test")
	RunnerInit(&cm)
	logContextTestRunner = LoggerContext{
		name:  "RUNNER_TEST",
		level: cm.GetLogLevel("general", 3)}

	fmt.Print("\n\nUno Test..\n\n")
	runner("uno")
	fmt.Print("\n\nDue Test..\n\n")
	runner("due")
	fmt.Print("\n\nTre Test..\n\n")
	runner("tre")

}
