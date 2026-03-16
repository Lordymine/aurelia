package cron

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Scheduler struct {
	store   Store
	runtime Runtime
	clock   Clock
	config  SchedulerConfig
}

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
	startedAt := now
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

func computeNextRun(expr string, after time.Time) (time.Time, error) {
	parts := strings.Fields(strings.TrimSpace(expr))
	if len(parts) != 5 {
		return time.Time{}, fmt.Errorf("unsupported cron expression: %s", expr)
	}

	minutePart := parts[0]
	hourPart := parts[1]

	switch {
	case strings.HasPrefix(minutePart, "*/") && hourPart == "*":
		interval, err := time.ParseDuration(strings.TrimPrefix(minutePart, "*/") + "m")
		if err != nil || interval <= 0 {
			return time.Time{}, fmt.Errorf("unsupported cron expression: %s", expr)
		}
		return after.UTC().Truncate(time.Minute).Add(interval), nil
	case isExactNumber(minutePart) && isExactNumber(hourPart):
		minute, _ := parseInt(minutePart)
		hour, _ := parseInt(hourPart)
		next := time.Date(after.Year(), after.Month(), after.Day(), hour, minute, 0, 0, time.UTC)
		if !next.After(after.UTC()) {
			next = next.Add(24 * time.Hour)
		}
		return next, nil
	default:
		return time.Time{}, fmt.Errorf("unsupported cron expression: %s", expr)
	}
}

func isExactNumber(v string) bool {
	_, err := parseInt(v)
	return err == nil
}

func parseInt(v string) (int, error) {
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, err
	}
	return n, nil
}
