package scheduler

import (
	"context"
	"github.com/lan143/metrology-master/internal/job"
	"go.uber.org/zap"
	"time"
)

type Scheduler struct {
	jobs   []job.Job
	jobCNs []chan struct{}

	log *zap.Logger
}

func NewScheduler(log *zap.Logger) *Scheduler {
	return &Scheduler{
		log: log,
	}
}

func (s *Scheduler) AddJob(job job.Job) {
	shutdownCh := make(chan struct{})
	s.jobCNs = append(s.jobCNs, shutdownCh)
	s.jobs = append(s.jobs, job)
}

func (s *Scheduler) Run() error {
	for i := range s.jobs {
		go s.executeJob(s.jobs[i], s.jobCNs[i])
	}

	return nil
}

func (s *Scheduler) Shutdown() error {
	for i := range s.jobCNs {
		s.jobCNs[i] <- struct{}{}
	}

	return nil
}

func (s *Scheduler) executeJob(job job.Job, shutdownCh chan struct{}) {
	ticker := time.NewTicker(1 * time.Minute)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		select {
		case <-shutdownCh:
			return
		case <-ticker.C:
			err := job.Execute(ctx)
			if err != nil {
				s.log.Error("execute job", zap.Error(err))
			}
		}
	}
}
