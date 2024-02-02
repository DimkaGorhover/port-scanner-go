package main

import (
	"context"
	"sync"
)

type Job func(ctx context.Context) error

type JobsExecutor interface {
	Submit(job Job)
	Shutdown()
}

type jobsExecutor struct {
	ctx      context.Context
	cancel   func()
	jobsChan chan Job
	wg       *sync.WaitGroup
}

func NewJobsExecutor(ctx context.Context, jobs int, parallelism int) JobsExecutor {
	ctx, cancel := context.WithCancel(ctx)
	var wg sync.WaitGroup
	je := &jobsExecutor{
		ctx:      ctx,
		cancel:   cancel,
		jobsChan: make(chan Job, jobs),
		wg:       &wg,
	}
	je.start(parallelism)
	return je
}

func (je *jobsExecutor) start(parallelism int) {
	je.wg.Add(parallelism)
	for i := 0; i < parallelism; i++ {
		go func() {
			defer je.wg.Done()
			for {
				select {
				case job := <-je.jobsChan:
					_ = job(je.ctx)
				case <-je.ctx.Done():
					return
				}
			}
		}()
	}
}

func (je *jobsExecutor) Submit(job Job) {
	je.jobsChan <- job
}

func (je *jobsExecutor) Shutdown() {
	defer close(je.jobsChan)
	je.cancel()
	je.wg.Wait()
}
