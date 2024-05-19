package service

import (
	"context"

	"go-service/internal/user/model"
)

type UserService interface {
	Load(ctx context.Context, id string) (*model.User, error)
	Create(ctx context.Context, user *model.User) (int64, error)
	Update(ctx context.Context, user *model.User) (int64, error)
	Patch(ctx context.Context, user map[string]interface{}) (int64, error)
	Delete(ctx context.Context, id string) (int64, error)
	Search(ctx context.Context, filter *model.UserFilter) ([]model.User, int64, error)
}
