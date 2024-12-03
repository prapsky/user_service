package user_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	params "github.com/prapsky/user_service/entity/params/user/find_all_param"
	qb "github.com/prapsky/user_service/internal/repository/query_builder/user"
)

func TestNewFindAllUsersQueryBuilder_WithAllOptions(t *testing.T) {
	const (
		expectedQuerySyntax = "SELECT id, name, phone_number, username, created_at FROM users " +
			"WHERE name ILIKE $1 AND username ILIKE $2 AND created_at BETWEEN $3 AND $4 ORDER BY created_at DESC LIMIT 25 OFFSET 5"
		username = "ronaldo"
		name     = "Cristiano Ronaldo"
		offset   = uint64(5)
		limit    = uint64(25)
		orderBy  = params.OrderByCreatedAtDesc
	)
	var (
		startDateTime = time.Now()
		endDateTime   = startDateTime.Add(time.Hour)
		timeRange     = params.TimeRange{Start: startDateTime, End: endDateTime}
	)

	param := params.NewFindAllUsersParam(params.WithName(name),
		params.WithUsername(username),
		params.WithFilterByCreatedTimeRange(timeRange),
		params.WithOffset(offset),
		params.WithOrderByCreatedAt(orderBy),
		params.WithLimit(limit))

	builder := qb.NewFindAllQueryBuilder(param)
	result := builder.Build()

	assert.Equal(t, expectedQuerySyntax, result.Syntax)
}
