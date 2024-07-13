package app

import (
	"context"
	"database/sql"

	"github.com/core-go/health"
	"github.com/core-go/log/zap"
	h "github.com/core-go/sql/health"

	"go-service/internal/user"
)

type ApplicationContext struct {
	Health *health.Handler
	User   user.UserTransport
}

func NewApp(ctx context.Context, cfg Config) (*ApplicationContext, error) {
	db, err := sql.Open(cfg.Sql.Driver, cfg.Sql.DataSourceName)
	if err != nil {
		return nil, err
	}
	logError := log.LogError

	userHandler, err := user.NewUserHandler(db, logError)
	if err != nil {
		return nil, err
	}

	sqlChecker := h.NewHealthChecker(db)
	healthHandler := health.NewHandler(sqlChecker)

	return &ApplicationContext{
		Health: healthHandler,
		User:   userHandler,
	}, nil
}
