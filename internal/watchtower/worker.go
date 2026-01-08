package watchtower

import (
	"context"
	"log/slog"
	"time"

	"watchtower/internal/logging"

	"github.com/go-co-op/gocron/v2"
)

type Workers struct {
	watchTower *Service
	cron       gocron.Scheduler
	logger     *slog.Logger
}

func NewWorkers(wt *Service) (*Workers, error) {
	logger := logging.FromContext(context.Background()).With("service", "workers")
	s, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}

	return &Workers{
		watchTower: wt,
		cron:       s,
		logger:     logger,
	}, nil
}

func (w *Workers) AddJobs() error {

	if _, err := w.cron.NewJob(gocron.DurationJob(time.Minute*2), gocron.NewTask(func() {
		w.logger.Debug("Running syncing orgs worker")

		if err := w.watchTower.SyncOrgs(); err != nil {
			w.logger.Error("Error syncing orgs", "error", err)
		}
	}), gocron.WithEventListeners()); err != nil {
		return err
	}

	if _, err := w.cron.NewJob(gocron.DurationJob(time.Minute*10), gocron.NewTask(func() {
		w.logger.Debug("Running remove old notifications worker")

		if err := w.watchTower.DeleteOldNotifications(); err != nil {
			w.logger.Error("Error syncing orgs", "error", err)
		}
	})); err != nil {
		return err
	}

	return nil
}

func (w *Workers) Start(ctx context.Context) {
	w.logger.Debug("Starting workers")

	w.cron.Start()
}

func (w *Workers) Stop() {
	w.logger.Debug("Stopping workers")

	if err := w.cron.StopJobs(); err != nil {
		w.logger.Error("Error stopping org sync worker", "error", err)
	}

	if err := w.cron.Shutdown(); err != nil {
		w.logger.Error("Error shutting down org sync worker", "error", err)
	}
}
