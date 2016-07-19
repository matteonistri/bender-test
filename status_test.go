package main

import "testing"

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
