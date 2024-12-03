package find_all_users_param

import (
	"time"
)

type OrderByCreatedAt uint8

const (
	UndefinedOrderByCreatedAt OrderByCreatedAt = iota
	OrderByCreatedAtAsc
	OrderByCreatedAtDesc
)

func (o OrderByCreatedAt) IsOrderByCreatedAt() bool {
	return o != UndefinedOrderByCreatedAt
}

type TimeRange struct {
	Start time.Time
	End   time.Time
}

func (u TimeRange) IsFilteredByCreatedAtTimeRange() bool {
	return !(u.Start.IsZero() && u.End.IsZero()) && (u.Start.Before(u.End))
}

type FindAllUsersParam struct {
	Name                     string
	Username                 string
	Offset                   uint64
	Limit                    uint64
	FilterByCreatedTimeRange TimeRange
	OrderByCreatedAt
}

type Options func(param *FindAllUsersParam)

func NewFindAllUsersParam(opts ...Options) FindAllUsersParam {
	params := FindAllUsersParam{
		Offset: 0,
		Limit:  20,
	}

	for _, opt := range opts {
		opt(&params)
	}

	return params
}

func WithName(name string) Options {
	return func(param *FindAllUsersParam) {
		param.Name = name
	}
}

func WithUsername(username string) Options {
	return func(param *FindAllUsersParam) {
		param.Username = username
	}
}

func WithFilterByCreatedTimeRange(timeRange TimeRange) Options {
	return func(param *FindAllUsersParam) {
		param.FilterByCreatedTimeRange = timeRange
	}
}

func WithOffset(offset uint64) Options {
	return func(param *FindAllUsersParam) {
		param.Offset = offset
	}
}

func WithLimit(limit uint64) Options {
	return func(param *FindAllUsersParam) {
		param.Limit = limit
	}
}

func WithOrderByCreatedAt(orderBy OrderByCreatedAt) Options {
	return func(param *FindAllUsersParam) {
		param.OrderByCreatedAt = orderBy
	}
}
