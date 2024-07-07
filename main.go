package main

import (
	"context"
	"github.com/core-go/config"
	"github.com/core-go/core/header"
	"github.com/core-go/core/random"
	srv "github.com/core-go/core/server"
	mid "github.com/core-go/log/middleware"
	"github.com/core-go/log/strings"
	"github.com/core-go/log/zap"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"go-service/internal/app"
)

func main() {
	var cfg app.Config
	err := config.Load(&cfg, "configs/config")
	if err != nil {
		panic(err)
	}
	r := mux.NewRouter()

	log.Initialize(cfg.Log)
	r.Use(mid.BuildContext)
	logger := mid.NewMaskLogger(cfg.MiddleWare.Request, Mask, Mask)
	// logger := mid.NewLogger()
	if log.IsInfoEnable() {
		r.Use(mid.Logger(cfg.MiddleWare, log.InfoFields, logger))
	}
	headerHandler := header.NewHeaderHandler(cfg.Response, GenerateId)
	r.Use(headerHandler.HandleHeader())
	r.Use(mid.Recover(log.PanicMsg))

	ctx := context.Background()
	err = app.Route(ctx, r, cfg)
	if err != nil {
		panic(err)
	}
	log.Info(ctx, srv.ServerInfo(cfg.Server))
	server := srv.CreateServer(cfg.Server, r)
	if err = server.ListenAndServe(); err != nil {
		log.Error(ctx, err.Error())
	}
}
func GenerateId() string {
	return random.Random(16)
}
func Mask(obj map[string]interface{}) {
	v, ok := obj["phone"]
	if ok {
		s, ok2 := v.(string)
		if ok2 && len(s) > 3 {
			obj["phone"] = strings.Mask(s, 0, 3, "*")
		}
	}
}
