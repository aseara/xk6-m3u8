package redis

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/redis", new(RedisExt))
}

type RedisExt struct{}

func (*RedisExt) NewClient(addr string, password string, db int) *redis.Client {
	if addr == "" {
		addr = "localhost:6379"
	}
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
}

func (*RedisExt) Set(client *redis.Client, key string, value string, expiration time.Duration) {
	d1 := []byte("hello\ngo\n")
	err := os.WriteFile("dat1", d1, 0644)

	time.Sleep(5 * time.Second)

	if err != nil {
		ReportError(err, "Failed to write local file")
	}
	err = client.Set(context.Background(), key, value, expiration).Err()
	if err != nil {
		ReportError(err, "Failed to set the specified key/value pair")
	}
}

func (*RedisExt) Get(client *redis.Client, key string) string {
	val, err := client.Get(context.Background(), key).Result()
	if err != nil {
		ReportError(err, "Failed to get the specified key")
	}
	return val
}

func (*RedisExt) Del(client *redis.Client, key string) {
	err := client.Del(context.Background(), key).Err()
	if err != nil {
		ReportError(err, "Failed to remove the specified key")
	}
}

func (*RedisExt) Do(client *redis.Client, args ...interface{}) (interface{}, error) {
	val, err := client.Do(context.Background(), args...).Result()
	if err != nil {
		if err == redis.Nil {
			return "", fmt.Errorf("key does not exist: %w", err)
		}
		return "", err
	}
	return val, nil
}
