package service_user_register

import (
	"context"
	"time"

	"github.com/prapsky/user_service/common/errors"
	"github.com/prapsky/user_service/common/logger/zerolog"
	"github.com/prapsky/user_service/common/utils"
	"github.com/prapsky/user_service/entity"
	auth_svc "github.com/prapsky/user_service/service/auth"
	user_svc "github.com/prapsky/user_service/service/user"
)

type RegisterUserUseCase interface {
	Do(ctx context.Context, input RegisterUserInput) (string, error)
}

type RegisterUserService struct {
	repo    user_svc.UserRepository
	authSvc auth_svc.AuthService
	log     zerolog.Zerolog
}

func NewRegisterUserService(repo user_svc.UserRepository,
	authSvc auth_svc.AuthService,
	log zerolog.Zerolog) RegisterUserService {
	return RegisterUserService{repo: repo,
		authSvc: authSvc,
		log:     log.WithServiceName("RegisterUserService")}
}

type RegisterUserInput struct {
	Name        string
	PhoneNumber string
	Username    string
	Password    string
}

func (svc RegisterUserService) Do(ctx context.Context, input RegisterUserInput) (string, error) {
	existingUserID, uerr := svc.repo.FindByUsername(ctx, input.Username)
	if uerr != nil {
		return "", uerr
	}

	if existingUserID != 0 {
		return "", errors.ErrUsernameAlreadyExists
	}

	user, uerr := EntityUser(input)
	if uerr != nil {
		return "", uerr
	}

	userID, err := svc.repo.Insert(ctx, user)
	if err != nil {
		svc.log.WarnfWithContext(ctx, err, "Error insert user")
		return "", err
	}

	user.ID = userID
	token, err := svc.authSvc.CreateToken(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

func EntityUser(input RegisterUserInput) (*entity.User, error) {
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &entity.User{
		Name:        input.Name,
		PhoneNumber: input.PhoneNumber,
		Username:    input.Username,
		Password:    hashedPassword,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}
