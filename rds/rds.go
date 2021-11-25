package rds

import (
	"context"
	"encoding/json"
	"github.com/0xJacky/Homework-api/settings"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"log"
	"time"
)

var rdb *redis.Client
var ctx = context.Background()

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     settings.RedisSettings.Addr,
		Password: settings.RedisSettings.Password,
		DB:       settings.RedisSettings.DB,
	})
}

func GetNoPrefix(key string) (string, error) {
	return rdb.Get(ctx, key).Result()
}

func Get(key string) (string, error) {
	return rdb.Get(ctx, settings.RedisSettings.Prefix+":"+key).Result()
}

func Keys(pattern string) (result []string, err error) {
	result, err = rdb.Keys(ctx, settings.RedisSettings.Prefix+":"+pattern).Result()
	return
}

func Set(key string, value interface{}, exp time.Duration) error {
	return rdb.Set(ctx, settings.RedisSettings.Prefix+":"+key, value, exp).Err()
}

func DelNoPrefix(key ...string) error {
	return rdb.Del(ctx, key...).Err()
}

func Del(key ...string) error {
	for i := range key {
		key[i] = settings.RedisSettings.Prefix + ":" + key[i]
	}
	return rdb.Del(ctx, key...).Err()
}

func Test() {
	log.Println("Testing Redis...")
	key := uuid.New().String()
	err := Set(settings.RedisSettings.Prefix, key, 100*time.Second)
	if err != nil {
		panic(err)
	}
	_, err = Get(settings.RedisSettings.Prefix)
	if err != nil {
		panic(err)
	}
	err = Del(key)
	if err != nil {
		panic(err)
	}
	log.Println("Redis ok")
}

func Publish(channel string, data gin.H) error {
	m, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = rdb.Publish(ctx, channel, string(m)).Err()
	return err
}

func Subscribe(channel string) *redis.PubSub {
	return rdb.Subscribe(ctx, channel)
}
