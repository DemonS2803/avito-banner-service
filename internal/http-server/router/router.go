package router

import (
	"avito-banner-service/internal/http-server/handlers/auth/token"
	"avito-banner-service/internal/http-server/handlers/url/banner"
	"avito-banner-service/internal/http-server/handlers/url/user-banner"
	"avito-banner-service/internal/repositories/postgres"
	"avito-banner-service/internal/repositories/redis"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Routes(redisClient *redis.Redis, db *postgres.Storage) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(token.NewTokenHandler())

	router.Get("/ping", user_banner.Ping())
	router.Get("/user_banner", user_banner.GetBannerById(redisClient, db))
	router.Get("/banner", banner.GetBannersFiltered(redisClient, db))
	router.Post("/banner", banner.CreateBanner(redisClient, db))
	router.Patch("/banner/{id}", banner.UpdateBanner(redisClient, db))
	router.Delete("/banner/{id}", banner.DeleteBanner(redisClient, db))

	return router
}
