package cron

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	robfigcron "github.com/robfig/cron/v3"
)

// Scheduler polls for due cron jobs and executes them.
type Scheduler struct {
	store   Store
	runtime Runtime
	clock   Clock
	config  SchedulerConfig
	running sync.Map // jobID → struct{} to prevent concurrent execution
}

// NewScheduler creates a cron scheduler.
func NewScheduler(store Store, runtime Runtime, clock Clock, config SchedulerConfig) (*Scheduler, error) {
	if store == nil {
		return nil, fmt.Errorf("cron store is required")
	}
	if runtime == nil {
		return nil, fmt.Errorf("cron runtime is required")
	}
	if clock == nil {
		clock = realClock{}
	}
	if config.PollInterval <= 0 {
		config.PollInterval = time.Minute
	}
	return &Scheduler{
		store:   store,
		runtime: runtime,
		clock:   clock,
		config:  config,
	}, nil
}

// RunDueJobs executes all jobs that are due.
func (s *Scheduler) RunDueJobs(ctx context.Context) (int, error) {
	now := s.clock.Now().UTC()
	jobs, err := s.store.ListDueJobs(ctx, now, 50)
	if err != nil {
		return 0, err
	}

	processed := 0
	for _, job := range jobs {
		if err := s.runSingleJob(ctx, now, job); err != nil {
			return processed, err
		}
		processed++
	}
	return processed, nil
}

// Start begins the scheduler polling loop.
func (s *Scheduler) Start(ctx context.Context) error {
	ticker := time.NewTicker(s.config.PollInterval)
	defer ticker.Stop()

	for {
		if _, err := s.RunDueJobs(ctx); err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}

func (s *Scheduler) runSingleJob(ctx context.Context, now time.Time, job CronJob) error {
	if _, loaded := s.running.LoadOrStore(job.ID, struct{}{}); loaded {
		return nil // already running
	}
	defer s.running.Delete(job.ID)

	startedAt := now
	log.Printf("cron.scheduler: executing job %s", job.ID)

	output, runErr := s.runtime.ExecuteJob(ctx, job)
	finishedAt := s.clock.Now().UTC()

	exec := CronExecution{
		ID:         uuid.NewString(),
		JobID:      job.ID,
		StartedAt:  startedAt,
		FinishedAt: &finishedAt,
	}

	if runErr != nil {
		exec.Status = "failed"
		exec.ErrorMessage = runErr.Error()
		job.LastStatus = "failed"
		job.LastError = runErr.Error()
	} else {
		exec.Status = "success"
		exec.OutputSummary = output
		job.LastStatus = "success"
		job.LastError = ""
	}

	job.LastRunAt = &finishedAt

	if job.ScheduleType == "once" {
		job.Active = false
		job.NextRunAt = nil
	} else if strings.EqualFold(job.ScheduleType, "cron") {
		nextRunAt, err := computeNextRun(job.CronExpr, now)
		if err != nil {
			return err
		}
		job.NextRunAt = &nextRunAt
	}

	if err := s.store.RecordExecution(ctx, exec); err != nil {
		return err
	}
	if err := s.store.UpdateJob(ctx, job); err != nil {
		return err
	}
	return nil
}

// computeNextRun calculates the next run time for a standard cron expression.
func computeNextRun(expr string, after time.Time) (time.Time, error) {
	parser := robfigcron.NewParser(robfigcron.Minute | robfigcron.Hour | robfigcron.Dom | robfigcron.Month | robfigcron.Dow)
	sched, err := parser.Parse(expr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid cron expression %q: %w", expr, err)
	}
	return sched.Next(after), nil
}
