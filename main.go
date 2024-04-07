package main

import (
	"avito-banner-service/internal/configs/postgres"
	"avito-banner-service/internal/configs/redis"
	"avito-banner-service/internal/http-server/handlers/auth/token"
	"avito-banner-service/internal/http-server/handlers/url/banner"
	"avito-banner-service/internal/http-server/handlers/url/user-banner"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(token.NewTokenHandler())

	router.Get("/user_banner", user_banner.GetBannerById(redisClient, db))
	router.Get("/banner", banner.GetBannersFiltered(redisClient, db))
	router.Post("/banner", banner.CreateBanner(redisClient, db))
	router.Patch("/banner/{id}", banner.UpdateBanner(redisClient, db))
	router.Delete("/banner/{id}", banner.DeleteBanner(redisClient, db))

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("shutdown server")

}
