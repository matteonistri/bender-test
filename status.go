package main

import (
	"errors"
	"time"
)

// StatusModule ...
type StatusModule struct {
	State     string
	Name      string
	Timestamp time.Time
	Jobs      map[string]Job
}

//StatusInterface ...
type StatusInterface interface {
	SetState(Job)
	GetState() (bool, int)
	GetJob(string) Job
	GetJobs(string) []Job
	GetRunningJob() (Job, error)
}

// Daemon states enum
const (
	DaemonIdle    = "idle"
	DaemonWorking = "working"
)

// SetState stores the provided Job into a map and updates the server idle
// status
func (s *StatusModule) SetState(job Job) {
	s.Jobs[job.UUID] = job
	if job.Status == JobWorking {
		s.State = DaemonWorking
		return
	}

	s.State = DaemonIdle
	return
}

// GetState returns the current idle status and the number of stored jobs
func (s *StatusModule) GetState() (string, int) {
	return s.State, len(s.Jobs)
}

// GetJob looks for a job with the specified uuid and returns its pointer
// or 'nil' if not found
func (s *StatusModule) GetJob(uuid string) Job {
	job := s.Jobs[uuid]
	return job
}

// GetJobs looks for jobs matching the specified name and returns a list
// of their poiners
func (s *StatusModule) GetJobs(name string) []Job {
	var jobs []Job
	for _, v := range s.Jobs {
		if v.Name == name {
			jobs = append(jobs, v)
		}
	}

	return jobs
}

// GetRunningJob returns a pointer to the currently running job or 'nil' if
// there is no running job
func (s *StatusModule) GetRunningJob() (Job, error) {
	for _, v := range s.Jobs {
		if v.Status == JobWorking {
			return v, nil
		}
	}

	return Job{}, errors.New("no running jobs (idle)")
}

//StatusModuleInit initializes the status module
func StatusModuleInit(sm *StatusModule, cm *ConfigModule) {
	sm.State = DaemonIdle
	sm.Name = cm.Get("status", "servername", "bender")
	sm.Jobs = make(map[string]Job)
	sm.Timestamp = time.Now()
}
