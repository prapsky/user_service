package user_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	qb "github.com/prapsky/user_service/internal/repository/query_builder/user"
)

func TestFindByIDQueryBuilder_Build(t *testing.T) {
	t.Run("return result with expected query syntax and arguments length", func(t *testing.T) {
		const expectedQuery = "SELECT name, phone_number, username " +
			"FROM users WHERE id = $1 LIMIT 1"
		id := uint64(7)

		builder := qb.NewFindByIDQueryBuilder(id).Build()
		assert.Equal(t, expectedQuery, builder.Syntax)
	})
}
