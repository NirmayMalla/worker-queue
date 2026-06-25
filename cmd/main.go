package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go-job-project/internal/job"
	"go-job-project/internal/worker"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// Create job manager
var jm = job.NewJobManager()

var handlers = map[string]job.Handler {
	"Send_email": 		job.EmailHandler{},
	"Process_file": 	job.FileHandler{},
	"Process_image": 	job.ImageHandler{},
}

// Create a worker pool
var p = worker.NewPool(5, handlers)

func main() {

	// Routes
	http.HandleFunc("/jobs", handleJobs)
	http.HandleFunc("/jobs/", handleJob)

	p.Start(jm)	

	quit := make(chan os.Signal, 1)

	server := &http.Server{
		Addr: ":8080",
	}

	// Run server in goroutine
	go func() {
		fmt.Printf("\n[+] Server running on http://localhost:8080 (http://localhost:8080/jobs)\n\n")
		server.ListenAndServe()
	}()

	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<- quit
	

	fmt.Printf("\n[!] Shutting down server...\n\n")
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)

	defer cancel()
	server.Shutdown(ctx)
	p.Stop()
}


// GET/POST all jobs
func handleJobs(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintln(w, "hello")

	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		
		jobs := jm.GetAll()		
			
		json.NewEncoder(w).Encode(jobs)
		return
	}

	if r.Method == http.MethodPost {

		var req struct{
			Type 		string	`json:"type"`
			Payload string	`json:"payload"`
		}

		json.NewDecoder(r.Body).Decode(&req)

		someJob := job.Job{
			ID: strconv.FormatInt(time.Now().UnixNano(), 10),
			Type: req.Type,
			Payload: req.Payload,
			Status: job.StatusPending,
			CreatedAt: time.Now(),
		}
		
		jm.AddJob(someJob)		// database

		p.Submit(someJob)	// add to queue (channel)

		// SEND RESPONSE
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(someJob)
	
		return
	}
}


// GET one job
func handleJob(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	
	id := strings.TrimPrefix(path, "/jobs/")

	if id == "" || strings.Contains(id, "/") {
		http.NotFound(w, r)
		return
	}

	job := jm.GetOne(id)
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(job)
}
