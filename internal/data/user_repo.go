package data

import (
	"context"

	"server-template/internal/biz"
	"server-template/internal/data/queries"
)

type userRepo struct {
	data *Data
}

func NewUserRepo(data *Data) biz.UserRepo {
	return &userRepo{
		data: data,
	}
}

func (u *userRepo) CreateUser(ctx context.Context, name string) (int64, error) {
	lastInsertId, err := u.data.WithWrite(ctx).CreateUser(ctx, name)
	if err != nil {
		return 0, err
	}
	return lastInsertId, nil
}

func (u *userRepo) CreateUserDetail(ctx context.Context, id int64, email string) (int64, error) {
	lastInsertId, err := u.data.WithWrite(ctx).CreateUserDetail(ctx, queries.CreateUserDetailParams{
		UserID: id,
		Email:  email,
	})
	if err != nil {
		return 0, err
	}
	return lastInsertId, nil
}

func (u *userRepo) UpdateUserInfo(ctx context.Context, name, country string) (int64, error) {
	rowsAffected, err := u.data.WithWrite(ctx).UpdateUser(ctx, queries.UpdateUserParams{
		Name:    name,
		Country: country,
	})
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}
