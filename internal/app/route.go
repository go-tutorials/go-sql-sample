package app

import (
	"context"

	"github.com/gorilla/mux"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	PATCH  = "PATCH"
	DELETE = "DELETE"
)

func Route(ctx context.Context, r *mux.Router, cfg Config) error {
	app, err := NewApp(ctx, cfg)
	if err != nil {
		return err
	}
	r.HandleFunc("/health", app.Health.Check).Methods(GET)

	user := "/users"
	r.HandleFunc(user+"/search", app.User.Search).Methods(GET, POST)
	r.HandleFunc(user+"/{id}", app.User.Load).Methods(GET)
	r.HandleFunc(user, app.User.Create).Methods(POST)
	r.HandleFunc(user+"/{id}", app.User.Update).Methods(PUT)
	r.HandleFunc(user+"/{id}", app.User.Patch).Methods(PATCH)
	r.HandleFunc(user+"/{id}", app.User.Delete).Methods(DELETE)

	return nil
}
