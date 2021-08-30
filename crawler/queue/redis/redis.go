package redis

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	zapp_core "github.com/zly-app/zapp/core"
	"go.uber.org/zap"

	"github.com/zly-app/service/crawler/config"
	"github.com/zly-app/service/crawler/core"
)

type RedisQueue struct {
	client redis.UniversalClient
}

func (r *RedisQueue) Put(queueName string, seed core.ISeed, front bool) error {
	data, err := seed.Encode()
	if err != nil {
		return fmt.Errorf("seed编码失败: %v", err)
	}
	if front {
		return r.client.LPush(context.Background(), queueName, data).Err()
	}
	return r.client.RPush(context.Background(), queueName, data).Err()
}

func (r *RedisQueue) Pop(queueName string, front bool) (core.ISeed, error) {
	panic("implement me")
}

func (r *RedisQueue) CheckQueueIsEmpty() (bool, error) {
	return true, nil
}

func (r *RedisQueue) Close() error {
	return r.client.Close()
}

func NewRedisQueue(app zapp_core.IApp) core.IQueue {
	conf := newRedisConfig()
	confKey := fmt.Sprintf("services.%s.queue.redis", config.NowServiceType)
	err := app.GetConfig().Parse(confKey, &conf)
	if err == nil {
		err = conf.Check()
	}
	if err != nil {
		app.Fatal("redis队列配置错误", zap.Error(err))
	}

	var client redis.UniversalClient
	if conf.IsCluster {
		client = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:        strings.Split(conf.Address, ","),
			Username:     conf.UserName,
			Password:     conf.Password,
			MinIdleConns: conf.MinIdleConns,
			PoolSize:     conf.PoolSize,
			ReadTimeout:  time.Duration(conf.ReadTimeout) * time.Millisecond,
			WriteTimeout: time.Duration(conf.WriteTimeout) * time.Millisecond,
			DialTimeout:  time.Duration(conf.DialTimeout) * time.Millisecond,
		})
	} else {
		client = redis.NewClient(&redis.Options{
			Addr:         conf.Address,
			Username:     conf.UserName,
			Password:     conf.Password,
			DB:           conf.DB,
			MinIdleConns: conf.MinIdleConns,
			PoolSize:     conf.PoolSize,
			ReadTimeout:  time.Duration(conf.ReadTimeout) * time.Millisecond,
			WriteTimeout: time.Duration(conf.WriteTimeout) * time.Millisecond,
			DialTimeout:  time.Duration(conf.DialTimeout) * time.Millisecond,
		})
	}
	return &RedisQueue{client: client}
}
