package redis

import (
	"avito-banner-service/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"log/slog"
	"time"
)

type Redis struct {
	client *redis.Client
}

func New() (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	fmt.Println("redis success configured")

	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("oh no. redis not send PONG (" + pong + ")")
		panic(err)
	}

	return &Redis{client: client}, nil
}

func PutBanner(client Redis, tagId int, featureId int, banner models.UserBanner) {
	ctx := context.Background()
	data, err := json.Marshal(banner)
	if err != nil {
		slog.Error("error while converting user-banner to json")
		panic(err)
	}
	err = client.client.Set(ctx, fmt.Sprintf("%d %d", tagId, featureId), data, 5*time.Minute).Err()
	if err != nil {
		slog.Error("error while saving user-banner to redis")
		panic(err)
	}
}

func GetBannerById(redisClient Redis, tagId int, featureId int, banner interface{}) bool {
	ctx := context.Background()
	req := redisClient.client.Get(ctx, fmt.Sprintf("%d %d", tagId, featureId))
	if err := req.Err(); err != nil {
		slog.Info("unable to GET data. error: %v", err)
		return false
	}
	res, err := req.Result()
	if err != nil {
		slog.Info("unable to GET data. error: %v", err)
		return false
	}
	json.Unmarshal([]byte(res), &banner)
	return true
}
