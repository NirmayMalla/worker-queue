package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go-job-project/internal/job"
	"go-job-project/internal/worker"
	"go-job-project/internal/config"
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

var conf = config.Setup()

var handlers = map[string]job.Handler {
	"Send_email": 		job.EmailHandler{},
	"Process_file": 	job.FileHandler{},
	"Process_image": 	job.ImageHandler{},
}

// Create a worker pool
var p = worker.NewPool(conf.WorkerCount, handlers, conf.QueueSize)

func main() {
	
	server := &http.Server{Addr: conf.Port}
	
	// Run server in goroutine
	go func() {
		fmt.Printf("\n[+] Server running on http://localhost:8080 (http://localhost:8080/jobs)\n\n")
		server.ListenAndServe()
	}()


	// Routes
	http.HandleFunc("/jobs", handleJobs)
	http.HandleFunc("/jobs/", handleJob)

	p.Start(jm)	

	quit := make(chan os.Signal, 1)

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
			
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(jobs)
		
		return
	}

	if r.Method == http.MethodPost {

		var req struct{
			Type 		string	`json:"type"`
			Payload string	`json:"payload"`
		}

		err := json.NewDecoder(r.Body).Decode(&req)

		// Malformed Data
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if req.Type == "" {
			http.Error(w, "A job type is required", http.StatusBadRequest)
			return
		}

		someJob := job.Job{
			ID: strconv.FormatInt(time.Now().UnixNano(), 10),
			Type: req.Type,
			Payload: req.Payload,
			Status: job.StatusPending,
			CreatedAt: time.Now(),
		}
		
		jm.AddJob(someJob)		// database [map]

		p.Submit(someJob)	// add to queue (channel)

		// SEND RESPONSE
		w.Header().Set("Content-Type", "application/json")
		
		w.WriteHeader(http.StatusCreated)	
		json.NewEncoder(w).Encode(someJob)
		
		return
	}
}


// GET one job
func handleJob(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	
	id := strings.TrimPrefix(path, "/jobs/")

	if id == "" || strings.Contains(id, "/") {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	
	// When the use wants to retrive their processed work
	job := jm.GetOne(id)

	// Might return a zero-value job {ID: ""}

	if job.ID == "" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(job)
}
