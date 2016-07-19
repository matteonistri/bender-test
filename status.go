package main

type Status struct {
	Idle bool
	Jobs map[string]Job
}

type StatusInterface interface {
	SetState(Job)
	GetState() (bool, int)
	GetJob(string) *Job
	GetJobs(string) *Job
	GetRunningJob() *Job
}

type StatusModule struct {
	Current Status
}

// SetState stores the provided Job into a map and updates the server idle
// status
func (s *StatusModule) SetState(job Job) {
	s.Current.Jobs[job.Uuid] = job
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

// GetJob looks for a job with the specified uuid and returns its pointer
// or 'nil' if not found
func (s *StatusModule) GetJob(uuid string) *Job {
	job := s.Current.Jobs[uuid]
	return &job
}

// GetJobs looks for jobs matching the specified name and returns a list
// of their poiners
func (s *StatusModule) GetJobs(name string) []*Job {
	var jobs []*Job
	for _, v := range s.Current.Jobs {
		if v.Name == name {
			jobs = append(jobs, &v)
		}
	}

	return jobs
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
