package watchtower

import (
	"context"
	"time"

	"watchtower/internal/logging"

	"github.com/go-co-op/gocron/v2"
)

type Workers struct {
	watchTower *Service
	cron       gocron.Scheduler
}

func NewWorkers(wt *Service) (*Workers, error) {
	s, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}

	return &Workers{
		watchTower: wt,
		cron:       s,
	}, nil
}

func (w *Workers) AddJobs() error {
	logger := logging.FromContext(context.Background()).With("service", "workers")

	if _, err := w.cron.NewJob(gocron.DurationJob(time.Minute*2), gocron.NewTask(func() {
		logger.Debug("Running syncing orgs worker")

		if err := w.watchTower.SyncOrgs(); err != nil {
			logger.Error("Error syncing orgs", "error", err)
		}
	})); err != nil {
		return err
	}

	if _, err := w.cron.NewJob(gocron.DurationJob(time.Minute*10), gocron.NewTask(func() {
		logger.Debug("Running remove old notifications worker")

		if err := w.watchTower.DeleteOldNotifications(); err != nil {
			logger.Error("Error syncing orgs", "error", err)
		}
	})); err != nil {
		return err
	}

	return nil
}

func (w *Workers) Start(ctx context.Context) {
	logger := logging.FromContext(ctx)
	logger.Debug("Starting workers")
	w.cron.Start()
}

func (w *Workers) Stop() {
	if err := w.cron.StopJobs(); err != nil {
		logging.FromContext(context.Background()).Error("Error stopping org sync worker", "error", err)
	}
}
