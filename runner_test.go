package main

import (
	"testing"
	"time"
)

func runner(name string, timeout time.Duration) (string, int) {
	job := &Job{}
	ret := job.Run(name, "f3628fa3-2b25-4dfd-ac2d-2c8a8613915c", []string{})
	if ret < 0 {
		return JobNotFound, 0
	}
	logChannel := job.Log()
	stateChannel := job.State()

	count := 0
	previousState := ""
	for {
		select {
		case m := <-logChannel:
			count += len(m)
		case s := <-stateChannel:
			if previousState != s {
				LogDeb(logContextTestRunner, "Receive [%v] state [%v]", name, s)
				previousState = s
			}
			if s != JobWorking {
				LogInf(logContextTestRunner, "%v", job)
				return s, count
			}
		case <-time.After(timeout * time.Second):
			LogDeb(logContextTestRunner, "Exec script [%v] Timeout! [%v]", name, timeout*time.Second)
			LogInf(logContextTestRunner, "%v", job)
			return JobTimeout, count
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
	ret, count := runner("Unknow", 60)
	if !(ret == JobNotFound && count == 0) {
		t.Errorf("Unknow should not found,0 [%v,%v].", ret, count)
	}
	LogInf(logContextTestRunner, "Foo Test..")
	ret, count = runner("foo", 5)
	if !(ret == JobFailed && count == 0) {
		t.Errorf("Foo script should fail,0 [%v,%v].", ret, count)
	}
	LogInf(logContextTestRunner, "Uno Test..")
	ret, count = runner("uno", 5)
	if !(ret == JobCompleted && count == 10) {
		t.Errorf("Uno script should be completed,10 [%v,%v].", ret, count)
	}
	LogInf(logContextTestRunner, "Due Test..")
	ret, count = runner("due", 10)
	if !(ret == JobFailed && count == 1000) {
		t.Errorf("Due script should fail,10 [%v,%v].", ret, count)
	}
	LogInf(logContextTestRunner, "Tre Test..")
	ret, count = runner("tre", 1)
	if !(ret == JobTimeout && count == 0) {
		t.Errorf("Tre script should be in timeout,0 [%v,%v].", ret, count)
	}

}
