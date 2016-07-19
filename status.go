package main

import "reflect"

type Status struct {
	Idle bool
	Jobs map[string]Job
}

type StatusInterface interface {
	SetState(*Job)
	GetState() (bool, int)
	GetJob(string, string) *Job
	GetRunningJob() *Job
}

type StatusModule struct {
	Current Status
}

// SetState stores the provided Job into a map and updates the server idle
// status
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

// GetJob looks for a job with the specified name and params and returns
// its pointer or 'nil' if not found
func (s *StatusModule) GetJob(name string, params string) *Job {
	for _, v := range s.Current.Jobs {
		if v.Name == name && v.Params == params {
			return &v
		}
	}

	return nil
}

// GetRunningJob returns a pointer to the currently running job or 'nil' if
// there is no running job
func (s *StatusModule) GetRunningJob() *Job {
	for _, v := range s.Current.Jobs {
		if v.Status == JOB_WORKING {
			return &v
		}
	}

	return nil
}
