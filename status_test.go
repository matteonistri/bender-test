package main

import (
	"reflect"
	"testing"

	"github.com/satori/go.uuid"
)

func TestRunningJob(t *testing.T) {
	sm := StatusModule{
		State: DaemonIdle,
		Jobs:  make(map[string]Job)}
	fakejob := Job{
		Status: JobWorking}
	sm.SetState(fakejob)

	actualB, actualI := sm.GetState()
	if actualB == DaemonIdle {
		t.Error("Expected DaemonWorking, got", actualB)
	}
	if actualI != 1 {
		t.Error("Expected 1, got", actualI)
	}
}

func TestFinishedJob(t *testing.T) {
	sm := StatusModule{
		State: DaemonIdle,
		Jobs:  make(map[string]Job)}
	fakejob := Job{
		Status: JobCompleted}
	sm.SetState(fakejob)

	actualB, actualI := sm.GetState()
	if actualB == DaemonWorking {
		t.Error("Expected DaemonIdle, got", actualB)
	}
	if actualI != 1 {
		t.Error("Expected 1, got", actualI)
	}
}

func TestEmptyJob(t *testing.T) {
	sm := StatusModule{
		State: DaemonIdle,
		Jobs:  make(map[string]Job)}
	fakejob := Job{}
	sm.SetState(fakejob)

	actualB, actualI := sm.GetState()
	if actualB == DaemonWorking {
		t.Error("Expected DaemonIdle, got", actualB)
	}
	if actualI != 1 {
		t.Error("Expected 1, got", actualI)
	}
}

func TestGetRunningJobValid(t *testing.T) {
	sm := StatusModule{
		State: DaemonIdle,
		Jobs:  make(map[string]Job)}
	fakejob := Job{
		Name:   "foo",
		Status: JobWorking}
	sm.SetState(fakejob)

	actualJ, err := sm.GetRunningJob()
	if !reflect.DeepEqual(actualJ, fakejob) {
		t.Error("Expected *Job, got", actualJ)
	}
	if err != nil {
		t.Error("Expected nil, got", err)
	}
}

func TestGetRunningJobIdle(t *testing.T) {
	sm := StatusModule{
		State: DaemonIdle,
		Jobs:  make(map[string]Job)}
	fakejob := Job{
		Name:   "foo",
		Status: JobCompleted}
	sm.SetState(fakejob)

	_, err := sm.GetRunningJob()
	if err == nil {
		t.Error("Expected nil, got", err)
	}
}

func TestGetJobsValid(t *testing.T) {
	sm := StatusModule{
		State: DaemonIdle,
		Jobs:  make(map[string]Job)}

	fakejobA := Job{
		Name: "fakejob",
		UUID: uuid.NewV4().String()}
	fakejobB := Job{
		Name: "fakejob",
		UUID: uuid.NewV4().String()}
	fakejobC := Job{
		Name: "fakejob",
		UUID: uuid.NewV4().String()}

	sm.SetState(fakejobA)
	sm.SetState(fakejobB)
	sm.SetState(fakejobC)

	actualJobs := sm.GetJobs("fakejob")
	if len(actualJobs) != 3 {
		t.Error("Expected 3, got", len(actualJobs))
	}
}

func TestGetJobsEmpty(t *testing.T) {
	sm := StatusModule{
		State: DaemonIdle,
		Jobs:  make(map[string]Job)}

	fakejobA := Job{
		UUID: uuid.NewV4().String()}
	fakejobB := Job{
		UUID: uuid.NewV4().String()}
	fakejobC := Job{
		UUID: uuid.NewV4().String()}

	sm.SetState(fakejobA)
	sm.SetState(fakejobB)
	sm.SetState(fakejobC)

	actualJobs := sm.GetJobs("")
	if len(actualJobs) != 3 {
		t.Error("Expected 3, got", len(actualJobs))
	}
}

func TestGetJobsInvalid(t *testing.T) {
	sm := StatusModule{
		State: DaemonIdle,
		Jobs:  make(map[string]Job)}

	fakejobA := Job{
		Name: "foo",
		UUID: uuid.NewV4().String()}
	fakejobB := Job{
		Name: "bar",
		UUID: uuid.NewV4().String()}
	fakejobC := Job{
		Name: "biz",
		UUID: uuid.NewV4().String()}

	sm.SetState(fakejobA)
	sm.SetState(fakejobB)
	sm.SetState(fakejobC)

	actualJobs := sm.GetJobs("fake")
	if len(actualJobs) != 0 {
		t.Error("Expected 0, got", len(actualJobs))
	}
}

func TestGetJobValid(t *testing.T) {
	sm := StatusModule{
		State: DaemonIdle,
		Jobs:  make(map[string]Job)}
	fakejob := Job{
		Name: "foo",
		UUID: "bar"}

	sm.SetState(fakejob)
	actualJ := sm.GetJob("bar")
	if !reflect.DeepEqual(actualJ, fakejob) {
		t.Error("Expected", fakejob, ", got", actualJ)
	}
}

func TestGetJobEmpty(t *testing.T) {
	sm := StatusModule{
		State: DaemonIdle,
		Jobs:  make(map[string]Job)}

	fakejob := Job{}

	sm.SetState(fakejob)
	actualJ := sm.GetJob("")
	if !reflect.DeepEqual(actualJ, fakejob) {
		t.Error("Expected", fakejob, ", got", actualJ)
	}
}
