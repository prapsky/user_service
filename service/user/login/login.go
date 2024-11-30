package service_user_login

import (
	"context"

	"github.com/prapsky/user_service/common/errors"
	"github.com/prapsky/user_service/common/logger/zerolog"
	"github.com/prapsky/user_service/common/utils"
	auth_svc "github.com/prapsky/user_service/service/auth"
	user_svc "github.com/prapsky/user_service/service/user"
)

type LoginUserUseCase interface {
	Do(ctx context.Context, input LoginUserInput) (string, error)
}

type LoginUserService struct {
	repo    user_svc.UserRepository
	authSvc auth_svc.AuthService
	log     zerolog.Zerolog
}

func NewLoginUserService(repo user_svc.UserRepository,
	authSvc auth_svc.AuthService,
	log zerolog.Zerolog) LoginUserService {
	return LoginUserService{repo: repo,
		authSvc: authSvc,
		log:     log.WithServiceName("LoginUserService")}
}

type LoginUserInput struct {
	Username string
	Password string
}

func (svc LoginUserService) Do(ctx context.Context, input LoginUserInput) (string, error) {
	user, uerr := svc.repo.FindByUsername(ctx, input.Username)
	if uerr != nil {
		return "", uerr
	}

	if user == nil {
		return "", errors.ErrAccountNotRegistered
	}

	if cerr := utils.ComparePasswords(user.Password, input.Password); cerr != nil {
		return "", errors.ErrIncorrectPassword
	}

	token, err := svc.authSvc.CreateToken(user)
	if err != nil {
		return "", err
	}

	return token, nil
}
