package job

import (
	"maps"
	"sync"
	"time"
)

type Job struct {
	ID 				string		`json:"id"`
	Type 			string		`json:"type"`
	Payload		string		`json:"payload"`
	Status		Status		`json:"status"`
	CreatedAt time.Time	`json:"created_at"`	
}


type Status string

const (
	StatusPending 		Status = "Pending"
	StatusProcessing 	Status = "Processing"
	StatusDone				Status = "Done"
	StatusFailed 			Status = "Failed"
)


type Handler interface {
	Handle(j Job) error
}


type JobManager struct {
	mu		sync.RWMutex
	jobs	map[string]Job
}

// Return pointer to new JobManager
func NewJobManager() *JobManager {
	return &JobManager{
		Jobs: make(map[string]Job),
	}
}


// Methods of JobManager
func (jm *JobManager) AddJob(s Job) {
	jm.mu.Lock()
	defer jm.mu.Unlock()

	jm.jobs[s.ID] = s
}


func (jm *JobManager) UpdateStatus(id string, status Status) {
	jm.mu.Lock()
	defer jm.mu.Unlock()

	job := jm.jobs[id]
	job.Status = status
	jm.jobs[id] = job
}

func (jm *JobManager) GetOne(id string) Job {
	jm.mu.RLock()
	defer jm.mu.RUnlock()

	j := jm.jobs[id]
	return j
}

func (jm *JobManager) GetAll() map[string]Job {
	jm.mu.RLock()
	defer jm.mu.RUnlock()

	copyJobMap := maps.Clone(jm.jobs)
	return copyJobMap
}
