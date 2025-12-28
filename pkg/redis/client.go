package redis

import (
	"time"

	"github.com/HariPrasath-3/scheduler-service/pkg/config"
	"github.com/redis/go-redis/v9"
)

func NewClient(redisConfig *config.RedisConfig) redis.UniversalClient {
	opts := &redis.UniversalOptions{
		Addrs:        redisConfig.Hosts,
		Password:     redisConfig.Password,
		PoolSize:     int(redisConfig.PoolSize),
		MinIdleConns: int(redisConfig.MinIdleConns),

		DialTimeout:  time.Duration(redisConfig.DialTimeout) * time.Millisecond,
		ReadTimeout:  time.Duration(redisConfig.ReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(redisConfig.WriteTimeout) * time.Millisecond,

		ConnMaxIdleTime: time.Duration(redisConfig.IdleTimeout) * time.Millisecond,
	}

	// Standalone vs Cluster behavior
	if redisConfig.Cluster {
		opts.ReadOnly = redisConfig.ServeReadsFromSlaves
		opts.RouteRandomly = redisConfig.ServeReadsFromMasterAndSlaves
		opts.MaxRedirects = 8
	} else {
		// standalone redis expects exactly one address
		if len(opts.Addrs) > 0 {
			opts.Addrs = []string{opts.Addrs[0]}
		}
		opts.DB = 0
	}

	return redis.NewUniversalClient(opts)
}
