package cron

import (
	"database/sql"
	"fmt"
)

func scanCronJob(scanner interface{ Scan(dest ...any) error }) (*CronJob, error) {
	var job CronJob
	var runAt, lastRunAt, nextRunAt sql.NullTime
	var active int
	err := scanner.Scan(
		&job.ID,
		&job.OwnerUserID,
		&job.TargetChatID,
		&job.ScheduleType,
		&job.CronExpr,
		&runAt,
		&job.Prompt,
		&active,
		&lastRunAt,
		&nextRunAt,
		&job.LastStatus,
		&job.LastError,
		&job.CreatedAt,
		&job.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	job.Active = active == 1
	if runAt.Valid {
		ts := runAt.Time
		job.RunAt = &ts
	}
	if lastRunAt.Valid {
		ts := lastRunAt.Time
		job.LastRunAt = &ts
	}
	if nextRunAt.Valid {
		ts := nextRunAt.Time
		job.NextRunAt = &ts
	}
	return &job, nil
}

func scanCronJobs(rows *sql.Rows) ([]CronJob, error) {
	var jobs []CronJob
	for rows.Next() {
		job, err := scanCronJob(rows)
		if err != nil {
			return nil, fmt.Errorf("scan cron job row: %w", err)
		}
		jobs = append(jobs, *job)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("cron job rows: %w", err)
	}
	return jobs, nil
}

func boolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}
