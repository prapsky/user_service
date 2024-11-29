package user

import (
	"github.com/Masterminds/squirrel"
	consts "github.com/prapsky/user_service/internal/repository/consts/user"
	query_bulder "github.com/prapsky/user_service/internal/repository/query_builder"
)

type FindByUsernameQueryBuilder struct {
	username string
}

func NewFindByUsernameQueryBuilder(username string) FindByUsernameQueryBuilder {
	return FindByUsernameQueryBuilder{
		username: username,
	}
}

func (qb FindByUsernameQueryBuilder) Build() query_bulder.Result {
	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	builder := sq.Select(consts.IDColumn).
		From(consts.TableName).
		Where(squirrel.Eq{consts.UsernameColumn: qb.username}).Limit(1)

	sql, params, _ := builder.ToSql()

	return query_bulder.Result{Syntax: sql, Params: params}
}
