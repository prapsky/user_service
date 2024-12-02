package service_user

import (
	"context"

	"github.com/prapsky/user_service/entity"
)

type UserRepository interface {
	Insert(ctx context.Context, user *entity.User) (uint64, error)
	FindByUsername(ctx context.Context, username string) (*entity.User, error)
	FindByID(ctx context.Context, id uint64) (*entity.User, error)
}
