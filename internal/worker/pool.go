package worker

import (
	"go-job-project/internal/job"
	"sync"
	"fmt"
)

type Pool struct {
	count 		int
	queue 		chan job.Job 
	handlers 	map[string]job.Handler
	wg				sync.WaitGroup
	queueSize int
}

func NewPool(count int, handlers map[string]job.Handler, queueSize int) *Pool {
	return &Pool {
		count: count,
		queue: make(chan job.Job, queueSize),
		handlers : handlers,
	} 
}

func (p *Pool) Start(jm *job.JobManager) {
	for i := 0; i < p.count; i++ {
		workerId := i + 1	
		p.wg.Add(1)

		go func(id int) {
			defer p.wg.Done()

			for j := range p.queue {
				jm.UpdateStatus(j.ID, job.StatusProcessing)
				
				fmt.Printf("[ Worker %d ]: processing Job, ID: %s, Type: %s\n", id, j.ID, j.Type)
				handler, exists := p.handlers[j.Type]

				if !exists {
					jm.UpdateStatus(j.ID, job.StatusFailed)
					fmt.Printf("[ Worker %d ]: Error processing Job, ID: %s, Type: %s, No handler type\n", id, j.ID, j.Type)
					continue
				}
				
				error := handler.Handle(j)

				// If handling the job returns an error
				if error != nil {
					jm.UpdateStatus(j.ID, job.StatusFailed)
					continue
				}

				jm.UpdateStatus(j.ID, job.StatusDone)

				fmt.Printf("[ Worker %d ]: Job done, %s, %s\n",id, j.ID, j.Type)
			}	
		}(workerId)
	}
}

func (p *Pool) Submit(j job.Job) {
	p.queue <- j
}

func (p *Pool) Stop() {
	close(p.queue)
	p.wg.Wait()
	fmt.Println("All workers done")
}
