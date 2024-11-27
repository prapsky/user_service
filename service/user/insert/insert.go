package service_user_insert

import (
	"context"
	"time"

	"github.com/prapsky/user_service/common/errors"
	"github.com/prapsky/user_service/common/logger/zerolog"
	"github.com/prapsky/user_service/entity"
	auth_svc "github.com/prapsky/user_service/service/auth"
	user_svc "github.com/prapsky/user_service/service/user"
)

type InsertUserUseCase interface {
	Do(ctx context.Context, input InsertUserInput) (string, error)
}

type InsertUserService struct {
	repo    user_svc.UserRepository
	authSvc auth_svc.AuthService
	log     zerolog.Zerolog
}

func NewInsertUserService(repo user_svc.UserRepository,
	authSvc auth_svc.AuthService,
	log zerolog.Zerolog) InsertUserService {
	return InsertUserService{repo: repo,
		authSvc: authSvc,
		log:     log.WithServiceName("InsertUserService")}
}

type InsertUserInput struct {
	Username string
}

func (svc InsertUserService) Do(ctx context.Context, input InsertUserInput) (string, error) {
	existingUserID, uerr := svc.repo.FindByUsername(ctx, input.Username)
	if uerr != nil {
		return "", uerr
	}

	if existingUserID != 0 {
		return "", errors.ErrUsernameAlreadyExists
	}

	now := time.Now()
	user := &entity.User{
		Username:  input.Username,
		CreatedAt: now,
		UpdatedAt: now,
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
