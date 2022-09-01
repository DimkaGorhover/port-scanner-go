package main

import (
	"sync"
)

type Job func()

type JobsExecutor interface {
	Start()
	Submit(job Job)
	Wait()
}

type jobsExecutor struct {
	jobsChan    chan Job
	parallelism int
	wg          *sync.WaitGroup
}

func NewJobsExecutor(jobs int, parallelism int) JobsExecutor {
	var wg sync.WaitGroup
	return &jobsExecutor{
		jobsChan:    make(chan Job, jobs),
		parallelism: parallelism,
		wg:          &wg,
	}
}

func (je jobsExecutor) Start() {
	for i := 0; i < je.parallelism; i++ {
		go je.worker()
	}
}

func (je jobsExecutor) Submit(job Job) {
	je.wg.Add(1)
	je.jobsChan <- job
}

func (je jobsExecutor) worker() {
	for job := range je.jobsChan {
		job()
		je.wg.Done()
	}
}

func (je jobsExecutor) Wait() {
	close(je.jobsChan)
	je.wg.Wait()
}
