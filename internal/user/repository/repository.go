package repository

import (
	"context"

	"go-service/internal/user/model"
)

type UserRepository interface {
	All(ctx context.Context) ([]model.User, error)
	Load(ctx context.Context, id string) (*model.User, error)
	Create(ctx context.Context, user *model.User) (int64, error)
	Update(ctx context.Context, user *model.User) (int64, error)
	Patch(ctx context.Context, user map[string]interface{}) (int64, error)
	Delete(ctx context.Context, id string) (int64, error)
	Search(ctx context.Context, filter *model.UserFilter, limit int64, offset int64) ([]model.User, int64, error)
}
