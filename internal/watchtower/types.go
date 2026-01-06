package watchtower

import (
	"context"
	"watchtower/internal/organisations"
	"watchtower/internal/products"
)

type Service struct {
	ctx        context.Context
	orgSvc     *organisations.Service
	productSvc *products.Service
	ghClient   ghClient
}
