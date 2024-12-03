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

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/prapsky/user_service/common/logger/zerolog"
	params "github.com/prapsky/user_service/entity/params/user/find_all_param"
	consts "github.com/prapsky/user_service/internal/repository/consts/user"
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

func (suite *UserRepositoryTestSuite) TestUser_FindByUsername() {
	const expectedQuery = "SELECT id, name, phone_number, password_hash " +
		"FROM users WHERE username = $1 LIMIT 1"

	input := createValidUser()

	suite.Run("select query returns error", func() {
		defer suite.subReporter.Add(suite.T())()

		suite.mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
			WithArgs(input.Username).
			WillReturnError(sql.ErrConnDone)

		result, err := suite.repo.FindByUsername(context.TODO(), input.Username)

		suite.NotNil(err)
		suite.Nil(result)
	})

	suite.Run("scan row error", func() {
		defer suite.subReporter.Add(suite.T())()

		suite.mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
			WithArgs(input.Username).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "phone_number", "password_hash"}).
				AddRow(input.CreatedAt, input.Name, input.PhoneNumber, input.Password))

		result, err := suite.repo.FindByUsername(context.TODO(), input.Username)

		suite.NotNil(err)
		suite.Nil(result)
	})

	suite.Run("successfully query user with given username", func() {
		defer suite.subReporter.Add(suite.T())()

		suite.mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
			WithArgs(input.Username).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "phone_number", "password_hash"}).
				AddRow(input.ID, input.Name, input.PhoneNumber, input.Password))

		result, err := suite.repo.FindByUsername(context.TODO(), input.Username)

		suite.Nil(err)
		suite.NotNil(result)
	})
}

func (suite *UserRepositoryTestSuite) TestUser_FindByID() {
	const expectedQuery = "SELECT name, phone_number, username " +
		"FROM users WHERE id = $1 LIMIT 1"

	input := createValidUser()

	suite.Run("select query returns error", func() {
		defer suite.subReporter.Add(suite.T())()

		suite.mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
			WithArgs(input.ID).
			WillReturnError(sql.ErrConnDone)

		result, err := suite.repo.FindByID(context.TODO(), input.ID)

		suite.NotNil(err)
		suite.Nil(result)
	})

	suite.Run("successfully query user with given id", func() {
		defer suite.subReporter.Add(suite.T())()

		suite.mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
			WithArgs(input.ID).
			WillReturnRows(sqlmock.NewRows([]string{"name", "phone_number", "username"}).
				AddRow(input.Name, input.PhoneNumber, input.Username))

		result, err := suite.repo.FindByID(context.TODO(), input.ID)

		suite.Nil(err)
		suite.NotNil(result)
	})
}

func (suite *UserRepositoryTestSuite) TestUser_FindAll() {
	suite.Run("unexpected error raised while query", func() {
		const (
			selectQuery = `SELECT id, name, phone_number, username, created_at 
			FROM users WHERE username ILIKE $1 LIMIT 20`
			username = "ronaldo"
		)

		params := params.NewFindAllUsersParam(params.WithUsername(username))

		suite.mock.ExpectQuery(regexp.QuoteMeta(selectQuery)).WillReturnError(sql.ErrConnDone)

		result, err := suite.repo.FindAll(context.TODO(), params)
		suite.Error(err)
		suite.ErrorIs(sql.ErrConnDone, err)
		suite.NotNil(result)
	})

	suite.Run("success return more users but get unexpected error while scan", func() {
		const (
			selectQuery = `SELECT id, name, phone_number, username, created_at 
			FROM users WHERE username ILIKE $1 LIMIT 20`
			name        = "Cristiano Ronaldo"
			phoneNumber = "0856241"
			username    = "ronaldo"
		)

		params := params.NewFindAllUsersParam(params.WithUsername(username))

		suite.mock.ExpectQuery(regexp.QuoteMeta(selectQuery)).
			WillReturnRows(sqlmock.NewRows([]string{consts.IDColumn,
				consts.NameColumn,
				consts.PhoneNumberColumn,
				consts.UsernameColumn,
				consts.CreatedAtColumn}).
				AddRow("wrong_id",
					name,
					phoneNumber,
					username,
					time.Now()))

		result, err := suite.repo.FindAll(context.TODO(), params)
		suite.Nil(err)
		suite.Empty(result)
	})

	suite.Run("success return more users but get unexpected error while scan", func() {
		const (
			selectQuery = `SELECT id, name, phone_number, username, created_at 
			FROM users WHERE username ILIKE $1 LIMIT 20`
			name        = "Cristiano Ronaldo"
			phoneNumber = "0856241"
			username    = "ronaldo"
			id          = uint64(7)
		)

		params := params.NewFindAllUsersParam(params.WithUsername(username))

		suite.mock.ExpectQuery(regexp.QuoteMeta(selectQuery)).
			WillReturnRows(sqlmock.NewRows([]string{consts.IDColumn,
				consts.NameColumn,
				consts.PhoneNumberColumn,
				consts.UsernameColumn,
				consts.CreatedAtColumn}).
				AddRow(id,
					name,
					phoneNumber,
					username,
					time.Now()))

		result, err := suite.repo.FindAll(context.TODO(), params)
		suite.Nil(err)
		suite.NotEmpty(result)
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
