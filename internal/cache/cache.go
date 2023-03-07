package cache

import (
	"ccfiw/internal/config"
	"crypto/tls"
	"errors"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

var RedisClient redis.UniversalClient

func Setup(c config.Config) error {
	log.Info("cache: setting up Redis client")
	if len(c.Redis.Servers) == 0 {
		return errors.New("at least one redis server must be configured")
	}

	var tlsConfig *tls.Config
	if c.Redis.TLSEnabled {
		tlsConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	if c.Redis.Cluster {
		RedisClient = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:     c.Redis.Servers,
			PoolSize:  c.Redis.PoolSize,
			Password:  c.Redis.Password,
			TLSConfig: tlsConfig,
		})
	} else if c.Redis.MasterName != "" {
		RedisClient = redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:       c.Redis.MasterName,
			SentinelAddrs:    c.Redis.Servers,
			SentinelPassword: c.Redis.Password,
			DB:               c.Redis.Database,
			PoolSize:         c.Redis.PoolSize,
			Password:         c.Redis.Password,
			TLSConfig:        tlsConfig,
		})
	} else {
		RedisClient = redis.NewClient(&redis.Options{
			Addr:      c.Redis.Servers[0],
			DB:        c.Redis.Database,
			Password:  c.Redis.Password,
			PoolSize:  c.Redis.PoolSize,
			TLSConfig: tlsConfig,
		})
	}

	return nil
}
