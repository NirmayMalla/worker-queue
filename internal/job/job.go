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
	Handle(j Job)
}


type JobManager struct {
	Mu		sync.RWMutex
	Jobs	map[string]Job
}

// Return pointer to new JobManager
func NewJobManager() *JobManager {
	return &JobManager{
		Jobs: make(map[string]Job),
	}
}


// Methods of JobManager
func (jm *JobManager) AddJob(s Job) {
	jm.Mu.Lock()
	defer jm.Mu.Unlock()

	jm.Jobs[s.ID] = s
}


func (jm *JobManager) UpdateStatus(id string, status Status) {
	jm.Mu.Lock()
	defer jm.Mu.Unlock()

	job := jm.Jobs[id]
	job.Status = status
	jm.Jobs[id] = job
}

func (jm *JobManager) GetOne(id string) Job {
	jm.Mu.RLock()
	defer jm.Mu.RUnlock()

	j := jm.Jobs[id]
	return j
}

func (jm *JobManager) GetAll() map[string]Job {
	jm.Mu.RLock()
	defer jm.Mu.RUnlock()

	copyJobMap := maps.Clone(jm.Jobs)
	return copyJobMap
}
