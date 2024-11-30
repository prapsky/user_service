package service_user_login_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"github.com/prapsky/user_service/common/logger/zerolog"
	"github.com/prapsky/user_service/common/utils"
	"github.com/prapsky/user_service/entity"
	svc "github.com/prapsky/user_service/service/user/login"
	mock_auth_service "github.com/prapsky/user_service/test/mock/service/auth"
	mock_user_service "github.com/prapsky/user_service/test/mock/service/user"
)

type LoginUserServiceTestSuite struct {
	suite.Suite
	mockCtrl *gomock.Controller
	repo     *mock_user_service.MockUserRepository
	authSvc  *mock_auth_service.MockAuthService
	svc      svc.LoginUserUseCase
}

func TestLoginUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(LoginUserServiceTestSuite))
}

func (suite *LoginUserServiceTestSuite) BeforeTest(string, string) {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.repo = mock_user_service.NewMockUserRepository(suite.mockCtrl)
	suite.authSvc = mock_auth_service.NewMockAuthService(suite.mockCtrl)
	suite.svc = svc.NewLoginUserService(suite.repo, suite.authSvc, zerolog.NewZeroLog())
}

func (suite *LoginUserServiceTestSuite) AfterTest(string, string) {
	defer suite.mockCtrl.Finish()
}

func (suite *LoginUserServiceTestSuite) TestNewLoginUserService() {
	suite.Run("successfully create a new instance", func() {
		suite.NotNil(suite.svc)
	})
}

func (suite *LoginUserServiceTestSuite) TestLoginUserService_Do() {
	token := "secret-token"
	user := createValidUser()
	emptyUser := &entity.User{}
	input := LoginInput(user)

	suite.Run("error find user by username", func() {
		suite.repo.EXPECT().FindByUsername(gomock.Any(), gomock.Any()).Return(emptyUser, errors.New("error"))
		token, err := suite.svc.Do(context.TODO(), input)
		suite.Error(err)
		suite.Empty(token)
	})

	suite.Run("error login user", func() {
		suite.repo.EXPECT().FindByUsername(gomock.Any(), gomock.Any()).Return(emptyUser, nil)
		token, err := suite.svc.Do(context.TODO(), input)
		suite.Error(err)
		suite.Empty(token)
	})

	suite.Run("success login", func() {
		suite.repo.EXPECT().FindByUsername(gomock.Any(), gomock.Any()).Return(user, nil)
		suite.authSvc.EXPECT().CreateToken(gomock.Any()).Return(token, nil)
		token, err := suite.svc.Do(context.TODO(), input)
		suite.Nil(err)
		suite.NotEmpty(token)
	})
}

func createValidUser() *entity.User {
	name := "Cristiano Ronaldo"
	username := "ronaldo"
	phoneNumber := "0856241"
	password := "secretagent"
	now := time.Now()

	hashedPassword, _ := utils.HashPassword(password)

	return &entity.User{
		ID:          7,
		Name:        name,
		PhoneNumber: phoneNumber,
		Username:    username,
		Password:    hashedPassword,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func LoginInput(input *entity.User) svc.LoginUserInput {
	password := "secretagent"
	return svc.LoginUserInput{
		Username: input.Username,
		Password: password,
	}
}
