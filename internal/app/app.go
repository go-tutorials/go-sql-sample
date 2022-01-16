package app

import (
	"context"
	"github.com/core-go/health"
	"github.com/core-go/log"
	sv "github.com/core-go/service"
	v "github.com/core-go/service/v10"
	q "github.com/core-go/sql"
	"github.com/core-go/sql/query"
	_ "github.com/go-sql-driver/mysql"
	"reflect"

	. "go-service/internal/handler"
	. "go-service/internal/model"
	. "go-service/internal/repository"
	. "go-service/internal/service"
)

type ApplicationContext struct {
	Health *health.Handler
	User   UserHandler
}

func NewApp(ctx context.Context, conf Config) (*ApplicationContext, error) {
	db, err := q.OpenByConfig(conf.Sql)
	if err != nil {
		return nil, err
	}
	logError := log.ErrorMsg
	status := sv.InitializeStatus(conf.Status)
	action := sv.InitializeAction(conf.Action)
	validator := v.NewValidator()

	userType := reflect.TypeOf(User{})
	userQueryBuilder := query.NewBuilder(db, "users", userType)
	userSearchBuilder, err := q.NewSearchBuilder(db, userType, userQueryBuilder.BuildQuery)
	if err != nil {
		return nil, err
	}

	userRepository := NewUserRepository(db) // userRepository, err := NewUserClient(conf.Client, log.InfoFields)
	userService := NewUserService(userRepository)
	userHandler := NewUserHandler(userSearchBuilder.Search, userService, status, logError, validator.Validate, &action)

	sqlChecker := q.NewHealthChecker(db)
	healthHandler := health.NewHandler(sqlChecker)

	return &ApplicationContext{
		Health: healthHandler,
		User:   userHandler,
	}, nil
}
