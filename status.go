package main

import "reflect"

type Status struct {
	Idle bool
	Jobs map[string]Job
}

type StatusInterface interface {
	SetState(*Job)
	GetState() (bool, int)
}

type StatusModule struct {
	Current Status
}

// State stores the provided Job into a map and updates the server idle status
func (s *StatusModule) SetState(job *Job) {
	if reflect.DeepEqual(job, &Job{}) || job == nil {
		LogAppendLine("STATUS  error: empty job provided")
		return
	}

	s.Current.Jobs[job.Uuid] = *job
	if job.Status == JOB_WORKING {
		s.Current.Idle = false
		return
	}

	s.Current.Idle = true
	return
}

// GetState returns the current idle status and the number of stored jobs
func (s *StatusModule) GetState() (bool, int) {
	return s.Current.Idle, len(s.Current.Jobs)
}

func init() {

}
