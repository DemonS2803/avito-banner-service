package mock_cache

import (
	ourRedis "avito-banner-service/internal/repositories/redis"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
)

func New() (*ourRedis.Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6380",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	fmt.Println("redis success configured")

	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("oh no. test redis not send PONG (" + pong + ")")
		panic(err)
	}

	return &ourRedis.Redis{Client: client}, nil
}
