package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type UserRepo interface {
	CreateUser(ctx context.Context, name string) (int64, error)
	UpdateUserInfo(ctx context.Context, name, country string) (int64, error)
	CreateUserDetail(ctx context.Context, id int64, email string) (int64, error)
}

type UserBiz struct {
	tx   Transaction
	repo UserRepo
	Rdb  redis.UniversalClient
	log  *log.Helper
}

func NewUserBiz(
	tx Transaction, repo UserRepo,
	redisCli redis.UniversalClient, logger log.Logger,
) *UserBiz {
	return &UserBiz{
		tx:   tx,
		repo: repo,
		Rdb:  redisCli,
		log:  log.NewHelper(log.With(logger, "module", "biz/user")),
	}
}

func (u *UserBiz) CreateUser(ctx context.Context, name, email string) (int64, error) {
	var (
		id  int64
		err error
	)

	err = u.tx.InTx(ctx, func(ctx context.Context) error {
		id, err = u.repo.CreateUser(ctx, name)
		if err != nil {
			return err
		}
		_, err = u.repo.CreateUserDetail(ctx, id, email)
		return err
	})
	if err != nil {
		return 0, errors.Wrap(err, "user create fail")
	}

	return id, nil
}
