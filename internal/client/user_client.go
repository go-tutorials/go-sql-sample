package client

import (
	"context"
	"net/http"

	"go-service/internal/model"
	"go-service/pkg/client"
)

type UserClient struct {
	Client *http.Client
	Url    string
	Config *client.LogConfig
	Log    func(context.Context, string, map[string]interface{})
}

type ResultInfo struct {
	Status  int64          `mapstructure:"status" json:"status" gorm:"column:status" bson:"status" dynamodbav:"status" firestore:"status"`
	Errors  []ErrorMessage `mapstructure:"errors" json:"errors,omitempty" gorm:"column:errors" bson:"errors,omitempty" dynamodbav:"errors,omitempty" firestore:"errors,omitempty"`
	Message string         `mapstructure:"message" json:"message,omitempty" gorm:"column:message" bson:"message,omitempty" dynamodbav:"message,omitempty" firestore:"message,omitempty"`
}
type ErrorMessage struct {
	Field   string `mapstructure:"field" json:"field,omitempty" gorm:"column:field" bson:"field,omitempty" dynamodbav:"field,omitempty" firestore:"field,omitempty"`
	Code    string `mapstructure:"code" json:"code,omitempty" gorm:"column:code" bson:"code,omitempty" dynamodbav:"code,omitempty" firestore:"code,omitempty"`
	Param   string `mapstructure:"param" json:"param,omitempty" gorm:"column:param" bson:"param,omitempty" dynamodbav:"param,omitempty" firestore:"param,omitempty"`
	Message string `mapstructure:"message" json:"message,omitempty" gorm:"column:message" bson:"message,omitempty" dynamodbav:"message,omitempty" firestore:"message,omitempty"`
}

func NewUserClient(config client.ClientConfig, log func(context.Context, string, map[string]interface{})) (*UserClient, error) {
	c, _, conf, err := client.InitializeClient(config)
	if err != nil {
		return nil, err
	}
	return &UserClient{Client: c, Url: config.Endpoint.Url, Config: conf, Log: log}, nil
}

func (c *UserClient) Load(ctx context.Context, id string) (*model.User, error) {
	url := c.Url + "/" + id
	var user model.User
	err := client.Get(ctx, c.Client, url, &user, c.Config, c.Log)
	return &user, err
}

func (c *UserClient) Create(ctx context.Context, user *model.User) (int64, error) {
	var res ResultInfo
	err := client.Post(ctx, c.Client, c.Url, user, &res, c.Config, c.Log)
	return res.Status, err
}

func (c *UserClient) Update(ctx context.Context, user *model.User, id string) (int64, error) {
	url := c.Url + "/" + id
	var res ResultInfo
	err := client.Put(ctx, c.Client, url, user, &res, c.Config, c.Log)
	return res.Status, err
}

func (c *UserClient) Patch(ctx context.Context, user map[string]interface{}, id string) (int64, error) {
	url := c.Url + "/" + id
	var res ResultInfo
	err := client.Patch(ctx, c.Client, url, user, &res, c.Config, c.Log)
	return res.Status, err
}

func (c *UserClient) Delete(ctx context.Context, id string) (int64, error) {
	url := c.Url + "/" + id
	var res int64
	err := client.Delete(ctx, c.Client, url, &res, c.Config, c.Log)
	return res, err
}
