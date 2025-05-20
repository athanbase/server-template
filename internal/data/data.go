package data

import (
	"context"
	"crypto/tls"
	"database/sql"
	"time"

	"server-template/internal/biz"
	"server-template/internal/conf"
	"server-template/internal/data/queries"

	"github.com/go-kratos/kratos/v2/log"

	redis "github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/wire"
	"github.com/pkg/errors"
)

var ProviderSet = wire.NewSet(NewTransaction, NewData, NewDB, NewUserRepo, NewRedis)

type contextTxKey struct{}

// Data provides all functions to execute db queries and transactions
type Data struct {
	masterConn    *sql.DB
	slaveConn     *sql.DB
	masterQueries *queries.Queries
	slaveQueries  *queries.Queries
	log           *log.Helper
	*queries.Queries
}

func (s *Data) WithRead() queries.Querier {
	return s.slaveQueries
}

func (s *Data) WithWrite(ctx context.Context) queries.Querier {
	tx, ok := ctx.Value(contextTxKey{}).(*queries.Queries)
	if ok {
		return tx
	}
	return s.masterQueries
}

func (s *Data) InTx(ctx context.Context, fn func(ctx context.Context) error) (err error) {
	tx, err := s.masterConn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if pErr := recover(); pErr != nil {
			s.log.Errorf("panic in tx, err: %+v", pErr)
			err = errors.Errorf("panic in tx, err: %+v", pErr)
			if rbErr := tx.Rollback(); rbErr != nil {
				s.log.Errorf("panic tx rollback failed: %+v", rbErr)
			}
			return
		}

		if err != nil {
			s.log.Warnf("tx got err, rolling back, err: %+v", err)
			if rbErr := tx.Rollback(); rbErr != nil {
				s.log.Errorf("tx rollback failed: %+v", rbErr)
			}
		}
	}()

	q := queries.New(tx)
	err = fn(context.WithValue(ctx, contextTxKey{}, q))
	if err != nil {
		return err // defer func will rollback
	}

	if err = tx.Commit(); err != nil {
		s.log.Warnf("tx commit err: %+v", err)
		return err // rollback in defer func
	}

	return nil
}

// NewTransaction .
func NewTransaction(d *Data) biz.Transaction {
	return d
}

// NewData creates a new store
func NewData(db *DB, logger log.Logger) (*Data, func(), error) {
	log := log.NewHelper(log.With(logger, "module", "data"))
	log.Info("init store")
	return &Data{
		masterConn:    db.master,
		masterQueries: queries.New(db.master),
		slaveConn:     db.slave,
		slaveQueries:  queries.New(db.slave),
		Queries:       queries.New(db.slave), // default use slave
		log:           log,
	}, func() {}, nil
}

type DB struct {
	master *sql.DB
	slave  *sql.DB
}

func NewDB(cfg *conf.DB, logger log.Logger) (*DB, func(), error) {
	log := log.NewHelper(log.With(logger, "module", "db"))
	log.Info("init db")

	masterCfg := cfg.Master
	master, err := sql.Open(masterCfg.Driver, masterCfg.Dsn)
	if err != nil {
		err = errors.Wrapf(err, "open master db failed, driver: %s", masterCfg.Driver)
		return nil, nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	err = master.PingContext(ctx)
	cancel()
	if err != nil {
		err = errors.Wrap(err, "master ping db failed")
		return nil, nil, err
	}

	master.SetMaxOpenConns(int(masterCfg.MaxOpenConn))
	master.SetMaxIdleConns(int(masterCfg.MaxIdleConn))
	master.SetConnMaxLifetime(
		time.Duration(masterCfg.MaxLifetimeConn) * time.Second,
	)

	slaveCfg := cfg.Slave
	slave, err := sql.Open(slaveCfg.Driver, slaveCfg.Dsn)
	if err != nil {
		err = errors.Wrapf(err, "open slave db failed, driver: %s", slaveCfg.Driver)
		return nil, nil, err
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*30)
	err = slave.PingContext(ctx)
	cancel()
	if err != nil {
		err = errors.Wrap(err, "slave ping db failed")
		return nil, nil, err
	}

	slave.SetMaxOpenConns(int(slaveCfg.MaxOpenConn))
	slave.SetMaxIdleConns(int(slaveCfg.MaxIdleConn))
	slave.SetConnMaxLifetime(
		time.Duration(slaveCfg.MaxLifetimeConn) * time.Second,
	)

	cleanup := func() {
		log.Info("closing the db connections")
		if err := master.Close(); err != nil {
			log.Errorf("close master conn failed, err: %+v", err)
		}

		if err := slave.Close(); err != nil {
			log.Errorf("close slave conn failed,err: %+v", err)
		}
	}
	return &DB{
		master: master,
		slave:  slave,
	}, cleanup, nil
}

func NewRedis(cfg *conf.Redis) (redis.UniversalClient, func(), error) {
	// init redis
	var tlsCfg *tls.Config
	if cfg.IsEnableTls {
		tlsCfg = &tls.Config{}
	}

	cli := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:      cfg.Addrs,
		Password:   cfg.Passwd,
		Username:   cfg.Username,
		MasterName: cfg.MasterName,
		TLSConfig:  tlsCfg,
		DB:         int(cfg.Db),
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(10))
	_, err := cli.Ping(ctx).Result()
	cancel()
	if err != nil {
		err = errors.Wrap(err, "ping redis failed")
		return nil, nil, err
	}

	log.Info("redis initialized")
	cleanup := func() {
		if err := cli.Close(); err != nil {
			log.Errorf("close redis conn failed, err: %+v", err)
		}
	}
	return cli, cleanup, nil
}
