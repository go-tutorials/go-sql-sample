package app

import (
	"context"
	. "github.com/core-go/service"
	"github.com/gorilla/mux"
)

func Route(r *mux.Router, ctx context.Context, conf Config) error {
	app, err := NewApp(ctx, conf)
	if err != nil {
		return err
	}

	r.HandleFunc("/health", app.Health.Check).Methods(GET)

	userPath := "/users"
	r.HandleFunc(userPath, app.User.Search).Methods(GET)
	r.HandleFunc(userPath+"/search", app.User.Search).Methods(GET, POST)
	r.HandleFunc(userPath+"/{id}", app.User.Load).Methods(GET)
	r.HandleFunc(userPath, app.User.Create).Methods(POST)
	r.HandleFunc(userPath+"/{id}", app.User.Update).Methods(PUT)
	r.HandleFunc(userPath+"/{id}", app.User.Patch).Methods(PATCH)
	r.HandleFunc(userPath+"/{id}", app.User.Delete).Methods(DELETE)

	return nil
}
