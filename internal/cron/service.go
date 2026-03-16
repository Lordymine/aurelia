package cron

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type Service struct {
	store Store
	clock Clock
}

func NewService(store Store, clock Clock) *Service {
	if clock == nil {
		clock = realClock{}
	}
	return &Service{store: store, clock: clock}
}

func (s *Service) CreateJob(ctx context.Context, job CronJob) (string, error) {
	if s.store == nil {
		return "", fmt.Errorf("cron store is required")
	}
	if strings.TrimSpace(job.Prompt) == "" {
		return "", fmt.Errorf("prompt is required")
	}
	if job.ID == "" {
		job.ID = uuid.NewString()
	}
	if job.LastStatus == "" {
		job.LastStatus = "idle"
	}
	job.Active = true
	now := s.clock.Now().UTC()
	if job.CreatedAt.IsZero() {
		job.CreatedAt = now
	}
	job.UpdatedAt = now

	switch strings.ToLower(strings.TrimSpace(job.ScheduleType)) {
	case "cron":
		if strings.TrimSpace(job.CronExpr) == "" {
			return "", fmt.Errorf("cron_expr is required for cron jobs")
		}
		nextRunAt, err := computeNextRun(job.CronExpr, now)
		if err != nil {
			return "", err
		}
		job.NextRunAt = &nextRunAt
	case "once":
		if job.RunAt == nil {
			return "", fmt.Errorf("run_at is required for once jobs")
		}
		job.NextRunAt = job.RunAt
	default:
		return "", fmt.Errorf("unsupported schedule_type %q", job.ScheduleType)
	}

	if err := s.store.CreateJob(ctx, job); err != nil {
		return "", err
	}
	return job.ID, nil
}

func (s *Service) ListJobs(ctx context.Context, chatID int64) ([]CronJob, error) {
	return s.store.ListJobsByChat(ctx, chatID)
}

func (s *Service) PauseJob(ctx context.Context, jobID string) error {
	job, err := s.store.GetJob(ctx, jobID)
	if err != nil {
		return err
	}
	if job == nil {
		return fmt.Errorf("cron job %s not found", jobID)
	}
	job.Active = false
	job.UpdatedAt = s.clock.Now().UTC()
	return s.store.UpdateJob(ctx, *job)
}

func (s *Service) ResumeJob(ctx context.Context, jobID string) error {
	job, err := s.store.GetJob(ctx, jobID)
	if err != nil {
		return err
	}
	if job == nil {
		return fmt.Errorf("cron job %s not found", jobID)
	}
	job.Active = true
	job.UpdatedAt = s.clock.Now().UTC()
	return s.store.UpdateJob(ctx, *job)
}

func (s *Service) DeleteJob(ctx context.Context, jobID string) error {
	return s.store.DeleteJob(ctx, jobID)
}
