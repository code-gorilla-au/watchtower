package main

import (
	"context"
	"watchtower/internal/watchtower"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx    context.Context
	worker *watchtower.OrgSyncWorker
}

// NewApp creates a new App application struct
func NewApp(worker *watchtower.OrgSyncWorker) *App {
	return &App{
		worker: worker,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.worker.Start(ctx)
}

func (a *App) shutdown(ctx context.Context) {
	a.worker.Stop()
}

func (a *App) OpenExternalURL(url string) {
	runtime.BrowserOpenURL(a.ctx, url)
}
