package service_user_insert_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"github.com/prapsky/user_service/common/logger/zerolog"
	svc "github.com/prapsky/user_service/service/user/insert"
	mock_auth_service "github.com/prapsky/user_service/test/mock/service/auth"
	mock_user_service "github.com/prapsky/user_service/test/mock/service/user"
)

type InsertUserServiceTestSuite struct {
	suite.Suite
	mockCtrl *gomock.Controller
	repo     *mock_user_service.MockUserRepository
	authSvc  *mock_auth_service.MockAuthService
	svc      svc.InsertUserUseCase
}

func TestInsertUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(InsertUserServiceTestSuite))
}

func (suite *InsertUserServiceTestSuite) BeforeTest(string, string) {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.repo = mock_user_service.NewMockUserRepository(suite.mockCtrl)
	suite.authSvc = mock_auth_service.NewMockAuthService(suite.mockCtrl)
	suite.svc = svc.NewInsertUserService(suite.repo, suite.authSvc, zerolog.NewZeroLog())
}

func (suite *InsertUserServiceTestSuite) AfterTest(string, string) {
	defer suite.mockCtrl.Finish()
}

func (suite *InsertUserServiceTestSuite) TestNewInsertUserService() {
	suite.Run("successfully create a new instance", func() {
		suite.NotNil(suite.svc)
	})
}

func (suite *InsertUserServiceTestSuite) TestInsertUserService_Do() {
	username := "ronaldo"
	token := "secret-token"
	userID := uint64(7)

	suite.Run("error find user by username", func() {
		suite.repo.EXPECT().FindByUsername(gomock.Any(), gomock.Any()).Return(uint64(0), errors.New("error"))
		token, err := suite.svc.Do(context.TODO(), svc.InsertUserInput{Username: username})
		suite.Error(err)
		suite.Empty(token)
	})

	suite.Run("username already exists", func() {
		suite.repo.EXPECT().FindByUsername(gomock.Any(), gomock.Any()).Return(userID, nil)
		token, err := suite.svc.Do(context.TODO(), svc.InsertUserInput{Username: username})
		suite.Error(err)
		suite.Empty(token)
	})

	suite.Run("error insert user", func() {
		suite.repo.EXPECT().FindByUsername(gomock.Any(), gomock.Any()).Return(uint64(0), nil)
		suite.repo.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(uint64(0), errors.New("some error"))
		token, err := suite.svc.Do(context.TODO(), svc.InsertUserInput{Username: username})
		suite.Error(err)
		suite.Empty(token)
	})

	suite.Run("success insert", func() {
		suite.repo.EXPECT().FindByUsername(gomock.Any(), gomock.Any()).Return(uint64(0), nil)
		suite.repo.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(userID, nil)
		suite.authSvc.EXPECT().CreateToken(gomock.Any()).Return(token, nil)
		token, err := suite.svc.Do(context.TODO(), svc.InsertUserInput{Username: username})
		suite.Nil(err)
		suite.NotEmpty(token)
	})
}
