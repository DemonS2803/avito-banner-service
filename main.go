package main

import (
	"avito-banner-service/internal/http-server/router"
	"avito-banner-service/internal/repositories/postgres"
	"avito-banner-service/internal/repositories/redis"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	fmt.Println("Hello World!")

	db, err := postgres.New()
	if err != nil {
		slog.Error("failed to init storage")
		os.Exit(1)
	}

	redisClient, err := redis.New()
	if err != nil {
		slog.Error("failed to init redis")
		os.Exit(1)
	}

	err = http.ListenAndServe(":8080", router.Routes(redisClient, db))
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("shutdown server")

}
