package app

import (
	"context"

	m "github.com/core-go/core"
	"github.com/gorilla/mux"
)

func Route(ctx context.Context, r *mux.Router, cfg Config) error {
	app, err := NewApp(ctx, cfg)
	if err != nil {
		return err
	}
	r.HandleFunc("/health", app.Health.Check).Methods(m.GET)

	user := "/users"
	r.HandleFunc(user+"/search", app.User.Search).Methods(m.GET, m.POST)
	r.HandleFunc(user+"/{id}", app.User.Load).Methods(m.GET)
	r.HandleFunc(user, app.User.Create).Methods(m.POST)
	r.HandleFunc(user+"/{id}", app.User.Update).Methods(m.PUT)
	r.HandleFunc(user+"/{id}", app.User.Patch).Methods(m.PATCH)
	r.HandleFunc(user+"/{id}", app.User.Delete).Methods(m.DELETE)

	return nil
}
