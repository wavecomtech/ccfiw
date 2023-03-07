package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

var ErrKeyNotFound = errors.New("key not found")

func GenKey(id string, mt string) string {
	now := time.Now().Truncate(time.Hour).Unix()
	return fmt.Sprintf("ccfiw::%s::%s::hr::%d", id, mt, now)
}

func SMAdd(pipe redis.Pipeliner, key string, vals []interface{}) error {
	if len(vals) == 0 {
		return nil
	}
	if _, err := pipe.SAdd(context.Background(), key, vals...).Result(); err != nil {
		log.Error("CacheSet SMAdd: ", err)
		return err
	}
	return nil
}

func GetMembersLength(key string) (int64, error) {
	b, err := RedisClient.SCard(context.Background(), key).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, ErrKeyNotFound
		}
		log.Error("GetMembers:", err)
		return 0, err
	}
	return b, nil
}

func SetExpiration(pipe redis.Pipeliner, key string) error {
	_, err := pipe.Expire(context.Background(), key, 48*time.Hour).Result()
	if err != nil {
		log.Error("SetExpiration:", err)
		return err
	}
	return nil
}
