package main

import (
	"context"
	"github.com/core-go/config"
	"github.com/core-go/core"
	"github.com/core-go/core/header"
	"github.com/core-go/core/random"
	"github.com/core-go/log/strings"
	"github.com/core-go/log/zap"
	mid "github.com/core-go/middleware"
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
	logger := mid.NewMaskLogger(MaskLog, MaskLog)
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
	log.Info(ctx, core.ServerInfo(cfg.Server))
	server := core.CreateServer(cfg.Server, r)
	if err = server.ListenAndServe(); err != nil {
		log.Error(ctx, err.Error())
	}
}
func GenerateId() string {
	return random.Random(16)
}
func MaskLog(name string, v interface{}) interface{}  {
	if name == "phone" {
		s, ok := v.(string)
		if ok {
			return strings.Mask(s, 0, 3, "*")
		}
	}
	return v
}
