package main

import (
	"fmt"
	"testing"
	"time"
)

func runner(name string) string {
	job := &Job{}
	ret := job.Run(name, "f3628fa3-2b25-4dfd-ac2d-2c8a8613915c", []string{})
	if ret < 0 {
		return JobNotFound
	}
	fmt.Println(ret)
	logChannel := *job.Log()
	stateChannel := *job.State()

	previousState := ""
	for {
		select {
		case m := <-logChannel:
			fmt.Println(m)
		case s := <-stateChannel:
			if previousState != s {
				LogDeb(logContextTestRunner, "Receive [%v] state [%v]", name, s)
				previousState = s
			}
			if s != JobWorking {
				return s
			}
		case <-time.After(60 * time.Second):
			LogDeb(logContextTestRunner, "Exec script [%v] Timeout!", name)
			return JobTimeout
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

	LogInf(logContextTestRunner, "Unknow Test..")
	ret := runner("Unknow")
	if ret != JobNotFound {
		t.Errorf("Unknow should not found [%v].", ret)
	}
	LogInf(logContextTestRunner, "Foo Test..")
	ret = runner("foo")
	if ret != JobFailed {
		t.Errorf("Foo script should fail [%v].", ret)
	}
	LogInf(logContextTestRunner, "Uno Test..")
	ret = runner("uno")
	if ret != JobFailed {
		t.Errorf("Uno script is not FAIL [%v].", ret)
	}
	LogInf(logContextTestRunner, "Due Test..")
	ret = runner("due")
	if ret != JobFailed {
		t.Errorf("Due script is not FAIL [%v].", ret)
	}
	LogInf(logContextTestRunner, "Tre Test..")
	ret = runner("tre")
	if ret != JobFailed {
		t.Errorf("Tre script is not FAIL [%v].", ret)
	}

}
