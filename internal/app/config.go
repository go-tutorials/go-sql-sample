package app

import (
	"github.com/core-go/core/header"
	"github.com/core-go/core/server"
	mid "github.com/core-go/log/middleware"
	"github.com/core-go/log/zap"
)

type Config struct {
	Server     server.ServerConf `mapstructure:"server"`
	Sql        SqlConfig         `mapstructure:"sql"`
	Log        log.Config        `mapstructure:"log"`
	Response   header.Config     `mapstructure:"response"`
	MiddleWare mid.LogConfig     `mapstructure:"middleware"`
}
type SqlConfig struct {
	DataSourceName string `yaml:"data_source_name" mapstructure:"data_source_name" json:"dataSourceName,omitempty" gorm:"column:datasourcename" bson:"dataSourceName,omitempty" dynamodbav:"dataSourceName,omitempty" firestore:"dataSourceName,omitempty"`
	Driver         string `yaml:"driver" mapstructure:"driver" json:"driver,omitempty" gorm:"column:driver" bson:"driver,omitempty" dynamodbav:"driver,omitempty" firestore:"driver,omitempty"`
}
