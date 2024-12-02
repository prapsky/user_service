package user

import (
	"github.com/Masterminds/squirrel"
	consts "github.com/prapsky/user_service/internal/repository/consts/user"
	query_bulder "github.com/prapsky/user_service/internal/repository/query_builder"
)

type FindByIDQueryBuilder struct {
	id uint64
}

func NewFindByIDQueryBuilder(id uint64) FindByIDQueryBuilder {
	return FindByIDQueryBuilder{
		id: id,
	}
}

func (qb FindByIDQueryBuilder) Build() query_bulder.Result {
	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	builder := sq.Select(consts.NameColumn, consts.PhoneNumberColumn, consts.UsernameColumn).
		From(consts.TableName).
		Where(squirrel.Eq{consts.IDColumn: qb.id}).Limit(1)

	sql, params, _ := builder.ToSql()

	return query_bulder.Result{Syntax: sql, Params: params}
}
