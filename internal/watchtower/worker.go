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
	if _, err := w.cron.NewJob(gocron.DurationJob(time.Minute*3), gocron.NewTask(func() {
		if err := w.watchTower.DeleteOldNotifications(); err != nil {
			logging.FromContext(context.Background()).Error("Error syncing orgs", "error", err)
		}
	})); err != nil {
		return err
	}

	if _, err := w.cron.NewJob(gocron.DurationJob(time.Minute*15), gocron.NewTask(func() {
		if err := w.watchTower.SyncOrgs(); err != nil {
			logging.FromContext(context.Background()).Error("Error syncing orgs", "error", err)
		}
	})); err != nil {
		return err
	}

	return nil
}

func (w *Workers) Start(ctx context.Context) {
	logger := logging.FromContext(ctx)
	logger.Debug("Starting org sync worker")
	w.cron.Start()
}

func (w *Workers) Stop() {
	if err := w.cron.StopJobs(); err != nil {
		logging.FromContext(context.Background()).Error("Error stopping org sync worker", "error", err)
	}
}
