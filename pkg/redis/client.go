package redis

import (
	"time"

	"github.com/HariPrasath-3/scheduler-service/pkg/config"
	"github.com/redis/go-redis/v9"
)

func NewClient(redisConfig *config.RedisConfig) *redis.ClusterClient {
	redisClusterClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:           redisConfig.Hosts,
		PoolSize:        int(redisConfig.PoolSize),
		MinIdleConns:    int(redisConfig.MinIdleConns),
		ReadOnly:        redisConfig.ServeReadsFromSlaves,
		RouteRandomly:   redisConfig.ServeReadsFromMasterAndSlaves,
		DialTimeout:     time.Duration(redisConfig.DialTimeout) * time.Millisecond,
		ReadTimeout:     time.Duration(redisConfig.ReadTimeout) * time.Millisecond,
		WriteTimeout:    time.Duration(redisConfig.WriteTimeout) * time.Millisecond,
		ConnMaxIdleTime: time.Duration(redisConfig.IdleTimeout) * time.Millisecond,
	})
	return redisClusterClient
}
