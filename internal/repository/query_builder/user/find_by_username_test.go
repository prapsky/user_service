package user_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	qb "github.com/prapsky/user_service/internal/repository/query_builder/user"
)

func TestFindByUsernameQueryBuilder_Build(t *testing.T) {
	t.Run("return result with expected query syntax and arguments length", func(t *testing.T) {
		const expectedQuery = "SELECT id, name, phone_number, password_hash " +
			"FROM users WHERE username = $1 LIMIT 1"
		username := "ronaldo"

		builder := qb.NewFindByUsernameQueryBuilder(username).Build()
		assert.Equal(t, expectedQuery, builder.Syntax)
	})
}
