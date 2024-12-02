package service_user_detail

import (
	"context"

	"github.com/prapsky/user_service/common/errors"
	"github.com/prapsky/user_service/common/logger/zerolog"
	"github.com/prapsky/user_service/entity"
	auth_svc "github.com/prapsky/user_service/service/auth"
	user_svc "github.com/prapsky/user_service/service/user"
)

type UserDetailUseCase interface {
	Do(ctx context.Context, input UserDetailInput) (*entity.User, error)
}

type UserDetailService struct {
	repo    user_svc.UserRepository
	authSvc auth_svc.AuthService
	log     zerolog.Zerolog
}

func NewUserDetailService(repo user_svc.UserRepository,
	authSvc auth_svc.AuthService,
	log zerolog.Zerolog) UserDetailService {
	return UserDetailService{repo: repo,
		authSvc: authSvc,
		log:     log.WithServiceName("UserDetailService")}
}

type UserDetailInput struct {
	Token string
}

func (svc UserDetailService) Do(ctx context.Context, input UserDetailInput) (*entity.User, error) {
	userID, verr := svc.authSvc.ValidateToken(input.Token)
	if verr != nil {
		return nil, errors.ErrInvalidToken
	}

	user, uerr := svc.repo.FindByID(ctx, userID)
	if uerr != nil {
		return nil, uerr
	}
	if user == nil {
		return nil, errors.ErrAccountNotRegistered
	}

	return user, nil
}
