package main

import (
	"reflect"
	"testing"

	"github.com/satori/go.uuid"
)

func TestRunningJob(t *testing.T) {
	sm := StatusModule{
		State: SERVER_IDLE,
		Jobs:  make(map[string]Job)}
	fakejob := Job{
		Status: JOB_WORKING}
	sm.SetState(fakejob)

	actualB, actualI := sm.GetState()
	if actualB == SERVER_IDLE {
		t.Error("Expected SERVER_WORKING, got", actualB)
	}
	if actualI != 1 {
		t.Error("Expected 1, got", actualI)
	}
}

func TestFinishedJob(t *testing.T) {
	sm := StatusModule{
		State: SERVER_IDLE,
		Jobs:  make(map[string]Job)}
	fakejob := Job{
		Status: JOB_COMPLETED}
	sm.SetState(fakejob)

	actualB, actualI := sm.GetState()
	if actualB == SERVER_WORKING {
		t.Error("Expected SERVER_IDLE, got", actualB)
	}
	if actualI != 1 {
		t.Error("Expected 1, got", actualI)
	}
}

func TestEmptyJob(t *testing.T) {
	sm := StatusModule{
		State: SERVER_IDLE,
		Jobs:  make(map[string]Job)}
	fakejob := Job{}
	sm.SetState(fakejob)

	actualB, actualI := sm.GetState()
	if actualB == SERVER_WORKING {
		t.Error("Expected SERVER_IDLE, got", actualB)
	}
	if actualI != 1 {
		t.Error("Expected 1, got", actualI)
	}
}

func TestGetRunningJobValid(t *testing.T) {
	sm := StatusModule{
		State: SERVER_IDLE,
		Jobs:  make(map[string]Job)}
	fakejob := Job{
		Name:   "foo",
		Status: JOB_WORKING}
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
		State: SERVER_IDLE,
		Jobs:  make(map[string]Job)}
	fakejob := Job{
		Name:   "foo",
		Status: JOB_COMPLETED}
	sm.SetState(fakejob)

	_, err := sm.GetRunningJob()
	if err == nil {
		t.Error("Expected nil, got", err)
	}
}

func TestGetJobsValid(t *testing.T) {
	sm := StatusModule{
		State: SERVER_IDLE,
		Jobs:  make(map[string]Job)}

	fakejobA := Job{
		Name: "fakejob",
		Uuid: uuid.NewV4().String()}
	fakejobB := Job{
		Name: "fakejob",
		Uuid: uuid.NewV4().String()}
	fakejobC := Job{
		Name: "fakejob",
		Uuid: uuid.NewV4().String()}

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
		State: SERVER_IDLE,
		Jobs:  make(map[string]Job)}

	fakejobA := Job{
		Uuid: uuid.NewV4().String()}
	fakejobB := Job{
		Uuid: uuid.NewV4().String()}
	fakejobC := Job{
		Uuid: uuid.NewV4().String()}

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
		State: SERVER_IDLE,
		Jobs:  make(map[string]Job)}

	fakejobA := Job{
		Name: "foo",
		Uuid: uuid.NewV4().String()}
	fakejobB := Job{
		Name: "bar",
		Uuid: uuid.NewV4().String()}
	fakejobC := Job{
		Name: "biz",
		Uuid: uuid.NewV4().String()}

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
		State: SERVER_IDLE,
		Jobs:  make(map[string]Job)}
	fakejob := Job{
		Name: "foo",
		Uuid: "bar"}

	sm.SetState(fakejob)
	actualJ := sm.GetJob("bar")
	if !reflect.DeepEqual(actualJ, fakejob) {
		t.Error("Expected", fakejob, ", got", actualJ)
	}
}

func TestGetJobEmpty(t *testing.T) {
	sm := StatusModule{
		State: SERVER_IDLE,
		Jobs:  make(map[string]Job)}

	fakejob := Job{}

	sm.SetState(fakejob)
	actualJ := sm.GetJob("")
	if !reflect.DeepEqual(actualJ, fakejob) {
		t.Error("Expected", fakejob, ", got", actualJ)
	}
}
