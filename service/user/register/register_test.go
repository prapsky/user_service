package service_user_register_test

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
	svc "github.com/prapsky/user_service/service/user/register"
	mock_auth_service "github.com/prapsky/user_service/test/mock/service/auth"
	mock_user_service "github.com/prapsky/user_service/test/mock/service/user"
)

type RegisterUserServiceTestSuite struct {
	suite.Suite
	mockCtrl *gomock.Controller
	repo     *mock_user_service.MockUserRepository
	authSvc  *mock_auth_service.MockAuthService
	svc      svc.RegisterUserUseCase
}

func TestRegisterUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(RegisterUserServiceTestSuite))
}

func (suite *RegisterUserServiceTestSuite) BeforeTest(string, string) {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.repo = mock_user_service.NewMockUserRepository(suite.mockCtrl)
	suite.authSvc = mock_auth_service.NewMockAuthService(suite.mockCtrl)
	suite.svc = svc.NewRegisterUserService(suite.repo, suite.authSvc, zerolog.NewZeroLog())
}

func (suite *RegisterUserServiceTestSuite) AfterTest(string, string) {
	defer suite.mockCtrl.Finish()
}

func (suite *RegisterUserServiceTestSuite) TestNewRegisterUserService() {
	suite.Run("successfully create a new instance", func() {
		suite.NotNil(suite.svc)
	})
}

func (suite *RegisterUserServiceTestSuite) TestRegisterUserService_Do() {
	user := createValidUser()
	emptyUser := &entity.User{}
	input := RegisterInput(user)

	suite.Run("error find user by username", func() {
		suite.repo.EXPECT().FindByUsername(gomock.Any(), gomock.Any()).Return(emptyUser, errors.New("error"))
		token, err := suite.svc.Do(context.TODO(), input)
		suite.Error(err)
		suite.Empty(token)
	})

	suite.Run("username already exists", func() {
		suite.repo.EXPECT().FindByUsername(gomock.Any(), gomock.Any()).Return(user, nil)
		token, err := suite.svc.Do(context.TODO(), input)
		suite.Error(err)
		suite.Empty(token)
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

func RegisterInput(input *entity.User) svc.RegisterUserInput {
	password := "secretagent"
	return svc.RegisterUserInput{
		Name:        input.Name,
		PhoneNumber: input.PhoneNumber,
		Username:    input.Username,
		Password:    password,
	}
}
