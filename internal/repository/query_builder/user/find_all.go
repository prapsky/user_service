package user

import (
	"fmt"

	"github.com/Masterminds/squirrel"
	param "github.com/prapsky/user_service/entity/params/user/find_all_param"
	consts "github.com/prapsky/user_service/internal/repository/consts/user"
	query_builder "github.com/prapsky/user_service/internal/repository/query_builder"
)

type FindAllQueryBuilder struct {
	params param.FindAllUsersParam
}

func NewFindAllQueryBuilder(params param.FindAllUsersParam) FindAllQueryBuilder {
	return FindAllQueryBuilder{
		params: params,
	}
}

func (f FindAllQueryBuilder) Build() query_builder.Result {
	const (
		defaultLimit = uint64(20)
	)
	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	builder := sq.Select(consts.IDColumn,
		consts.NameColumn,
		consts.PhoneNumberColumn,
		consts.UsernameColumn,
		consts.CreatedAtColumn).
		From(consts.TableName)

	builder = f.buildFilter(builder)

	if f.params.Offset > 0 {
		builder = builder.Offset(uint64(f.params.Offset))
	}

	if f.params.OrderByCreatedAt.IsOrderByCreatedAt() {
		if f.params.OrderByCreatedAt == param.OrderByCreatedAtAsc {
			builder = builder.OrderBy(fmt.Sprintf("%s ASC", consts.CreatedAtColumn))
		}

		if f.params.OrderByCreatedAt == param.OrderByCreatedAtDesc {
			builder = builder.OrderBy(fmt.Sprintf("%s DESC", consts.CreatedAtColumn))
		}
	}

	limit := uint64(f.params.Limit)
	if limit == uint64(0) {
		limit = defaultLimit
	}

	builder = builder.Limit(limit)

	sql, params := builder.MustSql()
	return query_builder.Result{Syntax: sql, Params: params}
}

func (f FindAllQueryBuilder) buildFilter(builder squirrel.SelectBuilder) squirrel.SelectBuilder {
	if f.params.Name != "" {
		name := fmt.Sprintf("%%%s%%", f.params.Name)
		builder = builder.Where(squirrel.Expr(fmt.Sprintf("%s ILIKE ?", consts.NameColumn), name))
	}

	if f.params.Username != "" {
		username := fmt.Sprintf("%%%s%%", f.params.Username)
		builder = builder.Where(squirrel.Expr(fmt.Sprintf("%s ILIKE ?", consts.UsernameColumn), username))
	}

	if f.params.FilterByCreatedTimeRange.IsFilteredByCreatedAtTimeRange() {
		builder = builder.Where(squirrel.Expr(fmt.Sprintf("%s BETWEEN ? AND ?", consts.CreatedAtColumn), f.params.FilterByCreatedTimeRange.Start, f.params.FilterByCreatedTimeRange.End))
	}

	return builder
}

func (f FindAllQueryBuilder) buildCount() query_builder.Result {
	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	builder := sq.Select(consts.Count).
		From(consts.TableName)
	builder = f.buildFilter(builder)

	sql, params := builder.MustSql()

	return query_builder.Result{
		Syntax: sql,
		Params: params,
	}
}
