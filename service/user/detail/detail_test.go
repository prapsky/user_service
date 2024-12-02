package service_user_detail_test

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
	svc "github.com/prapsky/user_service/service/user/detail"
	mock_auth_service "github.com/prapsky/user_service/test/mock/service/auth"
	mock_user_service "github.com/prapsky/user_service/test/mock/service/user"
)

type UserDetailServiceTestSuite struct {
	suite.Suite
	mockCtrl *gomock.Controller
	repo     *mock_user_service.MockUserRepository
	authSvc  *mock_auth_service.MockAuthService
	svc      svc.UserDetailUseCase
}

func TestUserDetailServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserDetailServiceTestSuite))
}

func (suite *UserDetailServiceTestSuite) BeforeTest(string, string) {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.repo = mock_user_service.NewMockUserRepository(suite.mockCtrl)
	suite.authSvc = mock_auth_service.NewMockAuthService(suite.mockCtrl)
	suite.svc = svc.NewUserDetailService(suite.repo, suite.authSvc, zerolog.NewZeroLog())
}

func (suite *UserDetailServiceTestSuite) AfterTest(string, string) {
	defer suite.mockCtrl.Finish()
}

func (suite *UserDetailServiceTestSuite) TestNewUserDetailService() {
	suite.Run("successfully create a new instance", func() {
		suite.NotNil(suite.svc)
	})
}

func (suite *UserDetailServiceTestSuite) TestUserDetailService_Do() {
	user := createValidUser()
	userID := user.ID
	input := DetailInput(user)

	suite.Run("error validate token", func() {
		suite.authSvc.EXPECT().ValidateToken(gomock.Any()).Return(uint64(0), errors.New("error"))
		userDetail, err := suite.svc.Do(context.TODO(), input)
		suite.Error(err)
		suite.Empty(userDetail)
	})

	suite.Run("error get user detail", func() {
		suite.authSvc.EXPECT().ValidateToken(gomock.Any()).Return(userID, nil)
		suite.repo.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
		userDetail, err := suite.svc.Do(context.TODO(), input)
		suite.Error(err)
		suite.Empty(userDetail)
	})

	suite.Run("success get detail user", func() {
		suite.authSvc.EXPECT().ValidateToken(gomock.Any()).Return(userID, nil)
		suite.repo.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(user, nil)
		userDetail, err := suite.svc.Do(context.TODO(), input)
		suite.Nil(err)
		suite.NotEmpty(userDetail)
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

func DetailInput(input *entity.User) svc.UserDetailInput {
	token := "secretagent"
	return svc.UserDetailInput{
		Token: token,
	}
}
