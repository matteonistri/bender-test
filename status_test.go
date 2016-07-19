package main

import (
	"reflect"
	"testing"
)

func TestRunningJob(t *testing.T) {
	status := Status{
		Idle: true,
		Jobs: make(map[string]Job)}
	sm := StatusModule{
		Current: status}
	fakejob := &Job{
		Status: JOB_WORKING}
	sm.SetState(fakejob)

	actualB, actualI := sm.GetState()
	if actualB {
		t.Error("Expected false, got", actualB)
	}
	if actualI != 1 {
		t.Error("Expected 1, got", actualI)
	}
}

func TestFinishedJob(t *testing.T) {
	status := Status{
		Idle: false,
		Jobs: make(map[string]Job)}
	sm := StatusModule{
		Current: status}
	fakejob := &Job{
		Status: JOB_COMPLETED}
	sm.SetState(fakejob)

	actualB, actualI := sm.GetState()
	if !actualB {
		t.Error("Expected true, got", actualB)
	}
	if actualI != 1 {
		t.Error("Expected 1, got", actualI)
	}
}

func TestEmptyJob(t *testing.T) {
	status := Status{
		Idle: true,
		Jobs: make(map[string]Job)}
	sm := StatusModule{
		Current: status}
	fakejob := &Job{}
	sm.SetState(fakejob)

	actualB, actualI := sm.GetState()
	if !actualB {
		t.Error("Expected true, got", actualB)
	}
	if actualI != 0 {
		t.Error("Expected 0, got", actualI)
	}
}

func TestNullJob(t *testing.T) {
	status := Status{
		Idle: true,
		Jobs: make(map[string]Job)}
	sm := StatusModule{
		Current: status}
	sm.SetState(nil)

	actualB, actualI := sm.GetState()
	if !actualB {
		t.Error("Expected true, got", actualB)
	}
	if actualI != 0 {
		t.Error("Expected 0, got", actualI)
	}
}

func TestGetJobValid(t *testing.T) {
	status := Status{
		Idle: true,
		Jobs: make(map[string]Job)}
	sm := StatusModule{
		Current: status}
	fakejob := &Job{
		Name:   "foo",
		Params: "--bar -v",
		Status: JOB_COMPLETED}
	sm.SetState(fakejob)
	actualJ := sm.GetJob("foo", "--bar -v")

	actualB, actualI := sm.GetState()
	if !actualB {
		t.Error("Expected true, got", actualB)
	}
	if actualI != 1 {
		t.Error("Expected 1, got", actualI)
	}
	if actualJ == nil {
		t.Error("Expected *Job, got", actualJ)
	}
	if !reflect.DeepEqual(fakejob, actualJ) {
		t.Error("Expected *Job, got", actualJ)
	}
}

func TestGetJobEmpty(t *testing.T) {
	status := Status{
		Idle: true,
		Jobs: make(map[string]Job)}
	sm := StatusModule{
		Current: status}
	fakejob := &Job{
		Status: JOB_COMPLETED}
	sm.SetState(fakejob)
	actualJ := sm.GetJob("", "")

	actualB, actualI := sm.GetState()
	if !actualB {
		t.Error("Expected true, got", actualB)
	}
	if actualI != 1 {
		t.Error("Expected 1, got", actualI)
	}
	if actualJ == nil {
		t.Error("Expected *Job, got", actualJ)
	}
	if !reflect.DeepEqual(fakejob, actualJ) {
		t.Error("Expected *Job, got", actualJ)
	}
}

func TestGetJobInvalid(t *testing.T) {
	status := Status{
		Idle: true,
		Jobs: make(map[string]Job)}
	sm := StatusModule{
		Current: status}
	fakejob := &Job{
		Name:   "foo",
		Status: JOB_COMPLETED}
	sm.SetState(fakejob)
	actualJ := sm.GetJob("bar", "")

	actualB, actualI := sm.GetState()
	if !actualB {
		t.Error("Expected true, got", actualB)
	}
	if actualI != 1 {
		t.Error("Expected 0, got", actualI)
	}
	if actualJ != nil {
		t.Error("Expected nil, got", actualJ)
	}
	if reflect.DeepEqual(fakejob, actualJ) {
		t.Error("Expected nil, got", actualJ)
	}
}
