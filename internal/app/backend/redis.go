package backend

import (
	"context"
	"emperror.dev/errors"
	"github.com/Sxtanna/chromatic_curator/internal/common"
	"github.com/Sxtanna/chromatic_curator/internal/system/backend"
	goredis "github.com/redis/go-redis/v9"
	"strconv"
)

const (
	configurationMissing = errors.Sentinel("redis configuration missing")
)

type RedisBackend struct {
	client *goredis.Client
}

func (r *RedisBackend) Init(config common.Configuration) error {
	redisConfiguration := common.FindConfiguration[backend.Config](config)
	if redisConfiguration == nil {
		return configurationMissing
	}

	r.client = goredis.NewClient(&goredis.Options{
		Addr:     redisConfiguration.Host + ":" + strconv.Itoa(redisConfiguration.Port),
		Password: "",
		DB:       0,
	})

	return nil
}

func (r *RedisBackend) Start() error {

	if _, err := r.client.Ping(context.Background()).Result(); err != nil {
		return errors.Wrap(err, "could not connect to redis database")
	}

	return common.ServiceStartedNormallyButDoesNotBlock
}

func (r *RedisBackend) Close(_ error) error {
	return r.client.Close()
}

func (r *RedisBackend) GetRole(ctx context.Context, user string) (string, error) {
	return r.client.HGet(ctx, "gamer_hoes:roles", user).Result()
}

func (r *RedisBackend) SetRole(ctx context.Context, user string, role string) error {
	return r.client.HSet(ctx, "gamer_hoes:roles", user, role).Err()
}
