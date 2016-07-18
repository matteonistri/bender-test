package main

import (
	"reflect"
	"time"
)

type JobStatus string

const (
	JOB_QUEUED     = "queued"
	JOB_NOT_FOUND  = "not found"
	JOB_QUEUE_FULL = "queue full"
	JOB_WORKING    = "working"
	JOB_FAILED     = "failed"
	JOB_COMPLETED  = "completed"
)

type Job struct {
	Script  string    `json:"script"`
	Path    string    `json:"path"`
	Args    []string  `json:"args"`
	Uuid    string    `json:"uuid"`
	Output  string    `json:"output"`
	Exit    string    `json:"exit"`
	Request time.Time `json:"request"`
	Start   time.Time `json:"start"`
	Finish  time.Time `json:"finish"`
	Status  JobStatus `json:"status"`
}

type Status struct {
	Idle bool           `json:"idle"`
	Jobs map[string]Job `json:"jobs"`
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
	if reflect.DeepEqual(job, &Job{}) {
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
