package app

import (
	"github.com/core-go/core"
	"github.com/core-go/log"
	mid "github.com/core-go/log/middleware"
	"github.com/core-go/sql"

	"go-service/pkg/client"
)

type Config struct {
	Server     core.ServerConf     `mapstructure:"server"`
	Sql        sql.Config          `mapstructure:"sql"`
	Client     client.ClientConfig `mapstructure:"client"`
	Log        log.Config          `mapstructure:"log"`
	MiddleWare mid.LogConfig       `mapstructure:"middleware"`
}
