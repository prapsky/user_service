package user_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/prapsky/user_service/entity"
	qb "github.com/prapsky/user_service/internal/repository/query_builder/user"
)

func TestInsertQueryBuilder_Build(t *testing.T) {
	t.Run("return result with expected query syntax and arguments length", func(t *testing.T) {
		var (
			name        = "Cristiano Ronaldo"
			password    = "goat"
			username    = "ronaldo"
			phoneNumber = "0856241"
		)

		const expectedQuery = "INSERT INTO users " +
			"(name,phone_number,username,password_hash,created_at,updated_at) " +
			"VALUES ($1,$2,$3,$4,$5,$6) " +
			"RETURNING id"
		now := time.Now()

		user := &entity.User{
			Name:        name,
			PhoneNumber: phoneNumber,
			Username:    username,
			Password:    password,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		builder := qb.NewInsertQueryBuilder(user).Build()
		assert.Equal(t, expectedQuery, builder.Syntax)
	})
}
