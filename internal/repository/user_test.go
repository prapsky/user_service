package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/prapsky/user_service/entity"
	"github.com/prapsky/user_service/internal/repository"
	"github.com/prapsky/user_service/test/helpers"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/prapsky/user_service/common/logger/zerolog"
	"github.com/stretchr/testify/suite"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	subReporter *helpers.SubReporter
	db          *sql.DB
	mock        sqlmock.Sqlmock
	repo        *repository.User
	log         zerolog.Zerolog
}

func TestUserRepository(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}

func (suite *UserRepositoryTestSuite) BeforeTest(string, string) {
	suite.subReporter = helpers.NewSubReporter(suite.T())

	db, mock, err := sqlmock.New()
	if err != nil {
		suite.FailNow("error opening a stub database connection: ", err)
	}

	suite.db = db
	suite.mock = mock
	suite.log = zerolog.NewZeroLog()
	suite.repo = repository.NewUser(db, suite.log)
}

func (suite *UserRepositoryTestSuite) AfterTest(string, string) {
	defer suite.db.Close()

	if err := suite.mock.ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (suite *UserRepositoryTestSuite) TestNewUserRepository() {
	suite.Run("successfully create a new instance of User Repository", func() {
		defer suite.subReporter.Add(suite.T())()
		suite.NotNil(suite.repo)
	})
}

func (suite *UserRepositoryTestSuite) TestUser_Insert() {
	const expectedQuery = "INSERT INTO users " +
		"(name,phone_number,username,password_hash,created_at,updated_at) " +
		"VALUES ($1,$2,$3,$4,$5,$6) RETURNING id"

	suite.Run("insert - error", func() {
		suite.mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WillReturnError(errors.New("some error"))
		userID, err := suite.repo.Insert(context.TODO(), &entity.User{})
		suite.Error(err)
		suite.Zero(userID)
	})

	suite.Run("insert - success", func() {
		suite.mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		userID, err := suite.repo.Insert(context.TODO(), createValidUser())
		suite.NoError(err)
		suite.NotZero(userID)
	})
}

func createValidUser() *entity.User {
	name := "Cristiano Ronaldo"
	username := "ronaldo"
	phoneNumber := "0856241"
	password := "secretagent"
	now := time.Now()

	return &entity.User{
		ID:          7,
		Name:        name,
		PhoneNumber: phoneNumber,
		Username:    username,
		Password:    password,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
