package watchtower

import (
	"context"
	"sync"
	"time"

	"watchtower/internal/logging"
)

type OrgSyncWorker struct {
	watchTower *Service
	wg         sync.WaitGroup
	stop       chan struct{}
}

func NewOrgSyncWorker(wt *Service) *OrgSyncWorker {
	return &OrgSyncWorker{
		watchTower: wt,
		wg:         sync.WaitGroup{},
		stop:       make(chan struct{}, 1),
	}
}

func (w *OrgSyncWorker) Start(ctx context.Context) {
	logger := logging.FromContext(ctx)
	logger.Debug("Starting org sync worker")

	w.wg.Add(1)

	go func() {
		for {
			select {
			case <-w.stop:
				logger.Debug("Stopping org sync worker")
				w.wg.Done()

				return
			default:
				if err := w.watchTower.SyncOrgs(); err != nil {
					logger.Error("Error syncing orgs", "error", err)
				}

				time.Sleep(time.Minute * 3)
			}
		}
	}()

	w.wg.Wait()
}

func (w *OrgSyncWorker) Stop() {
	w.stop <- struct{}{}
}
