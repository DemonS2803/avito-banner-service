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
	Client *redis.Client
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

	return &Redis{Client: client}, nil
}

func PutBanner(client Redis, tagId int, featureId int, banner models.UserBanner) {
	ctx := context.Background()
	data, err := json.Marshal(banner)
	if err != nil {
		slog.Error("error while converting user-banner to json")
		panic(err)
	}
	err = client.Client.Set(ctx, fmt.Sprintf("%d %d", tagId, featureId), data, 5*time.Minute).Err()
	if err != nil {
		slog.Error("error while saving user-banner to redis")
		panic(err)
	}
}

func PutEmptyBanner(client Redis, tagId int, featureId int) {
	ctx := context.Background()

	err := client.Client.Set(ctx, fmt.Sprintf("%d %d", tagId, featureId), "empty", 5*time.Minute).Err()
	if err != nil {
		slog.Error("error while saving user-banner to redis")
		panic(err)
	}
}

func IsEmptyBammer(client Redis, tagId int, featureId int) bool {
	ctx := context.Background()

	req := client.Client.Get(ctx, fmt.Sprintf("%d %d", tagId, featureId))
	if err := req.Err(); err != nil {
		slog.Info("unable to GET data. error: %v", err)
		return false
	}

	data, err := req.Result()
	if err != nil {
		slog.Info("unable to GET data. error: %v", err)
		return false
	}
	return data == "empty"

}

func GetBannerById(redisClient Redis, tagId int, featureId int, banner interface{}) bool {
	ctx := context.Background()
	req := redisClient.Client.Get(ctx, fmt.Sprintf("%d %d", tagId, featureId))
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

func GetBannerGroup(redisClient Redis, tagVal models.NilInt, featureVal models.NilInt, limit int, offset int) ([]models.UserBannerFilteredResponse, error) {
	ctx := context.Background()
	req := redisClient.Client.Get(ctx, fmt.Sprintf("group %d %d %d %d", tagVal.GetValue(), featureVal.GetValue(), limit, offset))
	if err := req.Err(); err != nil {
		slog.Info("unable to GET data. error: %v", err)
		return nil, err
	}
	data, _ := req.Result()
	var banners []models.UserBannerFilteredResponse
	err := json.Unmarshal([]byte(data), &banners)
	if err != nil {
		return nil, err
	}
	return banners, nil
}

func PutBannerGroup(redisClient Redis, tagVal models.NilInt, featureVal models.NilInt, banners []models.UserBannerFilteredResponse, limit int, offset int) {
	ctx := context.Background()
	data, err := json.Marshal(banners)
	if err != nil {
		slog.Info("unable to SET data. error: %v", err)
		return
	}
	req := redisClient.Client.Set(ctx, fmt.Sprintf("group %d %d %d %d", tagVal.GetValue(), featureVal.GetValue(), limit, offset), data, 5*time.Minute)
	if err := req.Err(); err != nil {
		slog.Info("unable to SET data. error: %v", err)
		return
	}
	return
}
