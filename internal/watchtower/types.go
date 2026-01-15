package watchtower

import (
	"context"
	"log/slog"
	"watchtower/internal/insights"
	"watchtower/internal/notifications"
	"watchtower/internal/organisations"
	"watchtower/internal/products"
)

type Service struct {
	ctx             context.Context
	orgSvc          *organisations.Service
	productSvc      *products.Service
	notificationSvc *notifications.Service
	insightSvc      *insights.Service
	ghClient        ghClient
}

type Config struct {
	Env      string
	AppDir   string
	LogLevel slog.Level
}
