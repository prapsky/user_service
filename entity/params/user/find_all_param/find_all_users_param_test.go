package find_all_users_param_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	params "github.com/prapsky/user_service/entity/params/user/find_all_param"
)

func TestNewFindAllUsersParam_WithNoOptions(t *testing.T) {
	param := params.NewFindAllUsersParam()

	assert.Equal(t, uint64(0), param.Offset)
	assert.Equal(t, uint64(20), param.Limit)
	assert.Equal(t, "", param.Name)
}

func TestNewFindAllUsersParam_WithOptions(t *testing.T) {
	const (
		limit    = uint64(10)
		offset   = uint64(20)
		name     = "Cristiano Ronaldo"
		username = "ronaldo"
	)

	var (
		startDateTime = time.Now()
		endDateTime   = startDateTime.Add(time.Hour)
		timeRange     = params.TimeRange{
			Start: startDateTime,
			End:   endDateTime,
		}
		orderByCreatedAt = params.OrderByCreatedAtDesc
	)

	param := params.NewFindAllUsersParam(
		params.WithName(name),
		params.WithUsername(username),
		params.WithFilterByCreatedTimeRange(timeRange),
		params.WithOffset(offset),
		params.WithLimit(limit),
		params.WithOrderByCreatedAt(orderByCreatedAt),
	)

	assert.Equal(t, name, param.Name)
	assert.Equal(t, username, param.Username)
	assert.Equal(t, true, param.FilterByCreatedTimeRange.IsFilteredByCreatedAtTimeRange())
	assert.Equal(t, offset, param.Offset)
	assert.Equal(t, limit, param.Limit)
	assert.Equal(t, true, param.OrderByCreatedAt.IsOrderByCreatedAt())
}
