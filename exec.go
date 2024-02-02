package main

import (
	"context"
	"runtime/pprof"
	"strconv"
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
	for i := 0; i < parallelism; i++ {
		labels := pprof.Labels(`worker`, strconv.Itoa(i))
		pprof.Do(je.ctx, labels, func(ctx context.Context) {
			je.runWorker(ctx)
		})
	}
}

func (je *jobsExecutor) runWorker(ctx context.Context) {
	je.wg.Add(1)
	go func() {
		defer je.wg.Done()
		for {
			select {
			case job, ok := <-je.jobsChan:
				if !ok {
					return
				}
				_ = job(ctx)
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (je *jobsExecutor) Submit(job Job) {
	je.jobsChan <- job
}

func (je *jobsExecutor) Shutdown() {
	defer je.cancel()
	close(je.jobsChan)
	je.wg.Wait()
}
