package app

import (
	"github.com/core-go/core/header"
	"github.com/core-go/core/server"
	mid "github.com/core-go/log/middleware"
	"github.com/core-go/log/zap"
	"github.com/core-go/sql"
)

type Config struct {
	Server     server.ServerConf `mapstructure:"server"`
	Sql        sql.Config        `mapstructure:"sql"`
	Log        log.Config        `mapstructure:"log"`
	Response   header.Config     `mapstructure:"response"`
	MiddleWare mid.LogConfig     `mapstructure:"middleware"`
}
