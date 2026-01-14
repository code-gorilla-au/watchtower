// Package watchtower provides the main application logic, including the service layer and background workers.
package watchtower

import (
	"context"
	"log/slog"
	"time"

	"watchtower/internal/logging"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Workers manages background jobs for syncing data and cleaning up old records.
type Workers struct {
	ctx        context.Context
	watchTower *Service
	cron       gocron.Scheduler
	logger     *slog.Logger
}

// NewWorkers initializes and returns a new Workers instance with the provided service.
func NewWorkers(wt *Service) (*Workers, error) {
	logger := logging.FromContext(context.Background()).With("service", "workers")
	s, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}

	return &Workers{
		ctx:        context.Background(),
		watchTower: wt,
		cron:       s,
		logger:     logger,
	}, nil
}

// AddJobs registers the background jobs to the scheduler.
func (w *Workers) AddJobs() error {

	if _, err := w.cron.NewJob(
		gocron.DurationJob(time.Minute*2),
		gocron.NewTask(w.jobSyncOrgs),
		gocron.WithEventListeners(gocron.AfterJobRuns(w.afterOrgSync)),
	); err != nil {
		return err
	}

	if _, err := w.cron.NewJob(
		gocron.DurationJob(time.Minute*10),
		gocron.NewTask(w.jobDeleteOldNotifications),
	); err != nil {
		return err
	}

	return nil
}

// Start begins the execution of scheduled jobs.
func (w *Workers) Start(ctx context.Context) {
	w.ctx = ctx

	w.logger.Debug("Starting workers")

	w.cron.Start()
}

// Stop halts the execution of scheduled jobs and shuts down the scheduler.
func (w *Workers) Stop() {
	w.logger.Debug("Stopping workers")

	if err := w.cron.StopJobs(); err != nil {
		w.logger.Error("Error stopping worker", "error", err)
	}

	if err := w.cron.Shutdown(); err != nil {
		w.logger.Error("Error shutting down worker", "error", err)
	}
}

func (w *Workers) jobSyncOrgs() {
	w.logger.Debug("Running syncing orgs worker")

	if err := w.watchTower.SyncOrgs(); err != nil {
		w.logger.Error("Error syncing orgs", "error", err)
	}
}

func (w *Workers) jobDeleteOldNotifications() {
	w.logger.Debug("Running remove old notifications worker")

	if err := w.watchTower.DeleteOldNotifications(); err != nil {
		w.logger.Error("Error syncing orgs", "error", err)
	}
}

func (w *Workers) afterOrgSync(jobID uuid.UUID, jobName string) {
	w.logger.Debug("Running notification worker")

	totalUnread, err := w.watchTower.CreateUnreadNotification()
	if err != nil {
		w.logger.Error("Error creating unread PR notification", "error", err)
	}

	if totalUnread > 0 {
		runtime.EventsEmit(w.ctx, "UNREAD_NOTIFICATIONS")
	}
}
