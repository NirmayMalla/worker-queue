# Worker Queue & Job Processing API

A concurrent job processing system built in Go.

The application exposes a REST API for creating and tracking background jobs while a configurable worker pool processes them asynchronously using goroutines and channels.

## Features
- REST API for submitting jobs
- Concurrent worker pool 
- Buffered job queue using Go channels
- Thread-safe in-memory job manager
- Job status tracking
- Configurable worker count
- Multiple job handlers using interfaces
- Graceful shutdown

## Concepts
- Goroutines
- Channels
- Worker pools
- Interfaces
- Mutexes
- WaitGroups
- Graceful shutdown
- REST API design
- JSON encoding decoding

## Project structure
```
cmd/
└─── main.go

internal
├─── config/
├─── job/
└─── worker/
```

##  Architecture
```
Client
   │
POST /jobs
   │
HTTP Handler
   │
Job Queue (channel)
   │
┌─────────────┐
│ Worker Pool │
└─────────────┘
   │
Handler
   │
Job Manager
```

A POST request creates a job and places it into the worker queue (buffered channel)

Workers read from the queue, update job status ("Processing")execute the appropriate handler and finally mark the job as "Done"


| Method | Endpoint |       Description        |
|--------|----------|--------------------------|
|  POST  |   /jobs  |      Create new job      |
|  GET   |   /jobs  |    Retrieve all jobs     |
|  GET   | jobs/id  |  Retrieve specific job   |


### Posting job:
#### POST ```/jobs```
```json
{
    "type": "Process_file"
    "payload": "example.txt"
}
```
### Queue full:
#### GET ```/jobs/{id}```
```json
{
    "id": "..."
    "status": "Pending"
}
```
### Before process completion:
#### GET ```/jobs/{id}```
```json
{
    "id": "..."
    "status": "Processing"
}
```
### After process time:
#### GET ```/jobs/{id}```
```json
{
    "id": "..."
    "status": "Done"
}
```

## Running

Clone the repository

cd go-job-project/cmd

go run main.go

The server starts on:

http://localhost:8080

