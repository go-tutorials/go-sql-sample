package app

import (
	"context"

	"github.com/core-go/health"
	"github.com/core-go/log/zap"
	q "github.com/core-go/sql"

	"go-service/internal/user"
)

type ApplicationContext struct {
	Health *health.Handler
	User   user.UserTransport
}

func NewApp(ctx context.Context, cfg Config) (*ApplicationContext, error) {
	db, err := q.OpenByConfig(cfg.Sql)
	if err != nil {
		return nil, err
	}
	logError := log.LogError

	userHandler, err := user.NewUserHandler(db, logError)
	if err != nil {
		return nil, err
	}

	sqlChecker := q.NewHealthChecker(db)
	healthHandler := health.NewHandler(sqlChecker)

	return &ApplicationContext{
		Health: healthHandler,
		User:   userHandler,
	}, nil
}
