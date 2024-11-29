package user

import (
	"github.com/Masterminds/squirrel"
	"github.com/prapsky/user_service/entity"
	consts "github.com/prapsky/user_service/internal/repository/consts/user"
	query_bulder "github.com/prapsky/user_service/internal/repository/query_builder"
)

type InsertQueryBuilder struct {
	user *entity.User
}

func NewInsertQueryBuilder(user *entity.User) InsertQueryBuilder {
	return InsertQueryBuilder{
		user: user,
	}
}

func (qb InsertQueryBuilder) Build() query_bulder.Result {
	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	builder := sq.Insert(consts.TableName).
		Columns(consts.NameColumn, consts.PhoneNumberColumn,
			consts.UsernameColumn, consts.PasswordHashColumn,
			consts.CreatedAtColumn, consts.UpdatedAtColumn).
		Values(qb.user.Name, qb.user.PhoneNumber,
			qb.user.Username, qb.user.Password,
			qb.user.CreatedAt, qb.user.UpdatedAt).
		SuffixExpr(squirrel.Expr("RETURNING id"))

	sql, params, _ := builder.ToSql()
	return query_bulder.Result{Syntax: sql, Params: params}
}
