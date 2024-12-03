package repository

import (
	"context"
	"database/sql"

	"github.com/prapsky/user_service/common/logger/zerolog"
	"github.com/prapsky/user_service/entity"
	params "github.com/prapsky/user_service/entity/params/user/find_all_param"
	queryBuilder "github.com/prapsky/user_service/internal/repository/query_builder/user"
)

type User struct {
	db  *sql.DB
	log zerolog.Zerolog
}

func NewUser(db *sql.DB, log zerolog.Zerolog) *User {
	return &User{
		db:  db,
		log: log.WithRepositoryName("User"),
	}
}

func (r *User) Insert(ctx context.Context, user *entity.User) (uint64, error) {
	query := queryBuilder.NewInsertQueryBuilder(user).Build()

	userID := uint64(0)
	if err := r.db.QueryRowContext(ctx, query.Syntax, query.Params...).Scan(&userID); err != nil {
		r.log.ErrorWithContext(ctx, err, "error exec insert query")
		return 0, err
	}

	return userID, nil
}

func (r *User) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	builder := queryBuilder.NewFindByUsernameQueryBuilder(username).Build()

	row := r.db.QueryRowContext(ctx, builder.Syntax, builder.Params...)
	var user entity.User
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.PhoneNumber,
		&user.Password,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (r *User) FindByID(ctx context.Context, id uint64) (*entity.User, error) {
	builder := queryBuilder.NewFindByIDQueryBuilder(id).Build()

	row := r.db.QueryRowContext(ctx, builder.Syntax, builder.Params...)
	var user entity.User
	err := row.Scan(
		&user.Name,
		&user.PhoneNumber,
		&user.Username,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (r *User) FindAll(ctx context.Context, params params.FindAllUsersParam) ([]*entity.User, error) {
	builder := queryBuilder.NewFindAllQueryBuilder(params)
	qb := builder.Build()

	rows, err := r.db.QueryContext(ctx, qb.Syntax, qb.Params...)
	if err != nil {
		return []*entity.User{}, err
	}
	defer rows.Close()

	users := []*entity.User{}
	for rows.Next() {
		user := entity.User{}
		err := rows.Scan(&user.ID,
			&user.Name,
			&user.PhoneNumber,
			&user.Username,
			&user.CreatedAt)

		if err != nil {
			r.log.ErrorfWithContext(ctx, err,
				"unexpected error while scan find all users ID#%d",
				user.ID)
			continue
		}

		users = append(users, &user)
	}

	return users, nil
}
