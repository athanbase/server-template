package sync

import (
	"context"
	"time"

	"server-template/pkg"
	"server-template/pkg/log"
	"server-template/pkg/retry"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

const (
	compareAndDel = "if redis.call('get', KEYS[1]) == ARGV[1] then return redis.call('del', KEYS[1]) else return 0 end"

	errWaitInternal = time.Millisecond * 200
)

func TryLock(ctxWithCancel context.Context, key, value string, expire, retry time.Duration, client redis.UniversalClient) {
	for {
		ok, err := Lock(ctxWithCancel, key, value, expire, client)
		if err != nil {
			log.Warnf("try lock %s error: %v", key, err)
			time.Sleep(errWaitInternal)
			continue
		}

		if !ok {
			log.Debugf("%s locked, try again later", key)
			time.Sleep(retry)
			continue
		}
		return
	}
}

func Lock(ctxWithCancel context.Context, key, value string, expire time.Duration, client redis.UniversalClient) (ok bool, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	ok, err = client.SetNX(ctx, key, value, expire).Result()
	cancel()
	if err != nil {
		err = errors.Wrapf(err, "set key: %s value: %s failed", key, value)
		return
	}

	// renew lock when obtain lock
	if ok {
		go func() {
			time.Sleep(expire - time.Second)
			for {
				select {
				case <-ctxWithCancel.Done():
					log.Debugf("got exit signal, pExpire task exit, key: %s value: %s", key, value)
					return
				default:
					ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
					ok, err1 := client.PExpire(ctx, key, expire).Result()
					cancel()
					if err1 != nil {
						log.Errorf("pExpire key: %s failed: %s", key, err1.Error())
						time.Sleep(time.Millisecond * 20)
						return
					}
					if !ok {
						log.Warnf("pExpire key: %s %s failed", key, expire)
						time.Sleep(time.Millisecond * 20)
						return
					}
					log.Debugf("pExpire key: %s %s success", key, expire)
					time.Sleep(expire - time.Second)
				}
			}
		}()
	}
	return
}

func Unlock(ctx context.Context, key, value string, cancelDaemonFunc context.CancelFunc, client redis.UniversalClient) (ok bool, err error) {
	cancelDaemonFunc()
	res, err := client.Eval(ctx, compareAndDel, []string{key}, value).Result()
	if err != nil {
		return
	}

	return res.(int64) == 1, nil
}

func Ping(client *redis.ClusterClient) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	_, err = client.Ping(ctx).Result()
	cancel()
	return
}

func TryLockWithKey(ctx context.Context, key string, client redis.UniversalClient) func() (bool, error) {
	ctx1, cancel := context.WithCancel(ctx)
	value := pkg.GenID()
	expire := time.Second * 30
	retry := time.Millisecond * 100
	TryLock(ctx1, key, value, expire, retry, client)

	return func() (bool, error) {
		ctx, cancel1 := context.WithTimeout(context.Background(), time.Second)
		defer cancel1()
		return Unlock(ctx, key, value, cancel, client)
	}
}

func LockWithKey(ctx context.Context, key string, client redis.UniversalClient) (isLocked bool, cancelFunc func() bool, err error) {
	ctx1, cancel := context.WithCancel(ctx)
	value := pkg.GenID()
	expire := time.Second * 30

	err = retry.BackoffRetry(func() error {
		isLocked, err = Lock(ctx1, key, value, expire, client)
		if err != nil {
			log.Warnf("lock %s failed, err: %+v, retry", key, err)
			return err
		}
		return nil
	})

	cancelFunc = func() bool {
		ctx, cancel1 := context.WithTimeout(context.Background(), time.Second)
		defer cancel1()
		ok, err := Unlock(ctx, key, value, cancel, client)
		if err != nil {
			log.Errorf("unlock failed, key: %s, value: %s", key, value)
		}
		return ok
	}

	return
}
