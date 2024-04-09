package banner

import (
	"avito-banner-service/internal/models"
	"avito-banner-service/internal/repositories/postgres"
	"avito-banner-service/internal/repositories/redis"
	resp "avito-banner-service/internal/utils/response"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"strconv"
)

func GetBanners(redisClient *redis.Redis, db *postgres.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("token") != "admin_token" {
			resp.Send403Error(w, r)
			return
		}

		tagId, err := strconv.Atoi(r.URL.Query().Get("tag_id"))
		if err != nil {
			slog.Error("error parsing param tag_id")
			resp.Send400Error(w, r)
			return
		}
		featureId, err := strconv.Atoi(r.URL.Query().Get("feature_id"))
		if err != nil {
			slog.Error("error parsing param feature_id")
			resp.Send400Error(w, r)
			return
		}

		var useLastRevision bool
		useLastRevisionStr := r.URL.Query().Get("use_last_revision")
		if useLastRevisionStr == "true" {
			useLastRevision = true
		} else {
			useLastRevision = false
		}

		fmt.Println("useLastRevision ")
		fmt.Println(useLastRevision)

		var banner models.UserBanner
		isCached := redis.GetBannerById(*redisClient, tagId, featureId, &banner)
		if useLastRevision || !isCached {
			banner, err = postgres.GetUserBannerByTagIdAndFeatureId(db, tagId, featureId)
			if err != nil {
				resp.Send404Error(w, r)
				return
			}
		}
		redis.PutBanner(*redisClient, tagId, featureId, banner)

		render.JSON(w, r, banner)
	}
}

func GetBannersFiltered(redisClient *redis.Redis, db *postgres.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("token") != "admin_token" {
			resp.Send403Error(w, r)
			return
		}
		tagVal := models.NilInt{}

		tagId, err := strconv.Atoi(r.URL.Query().Get("tag_id"))
		if err != nil {
			tagVal.Null = true
		} else {
			tagVal.Null = false
			tagVal.Value = tagId
		}

		featureVal := models.NilInt{}
		featureId, err := strconv.Atoi(r.URL.Query().Get("feature_id"))
		if err != nil {
			featureVal.Null = true
		} else {
			featureVal.Null = false
			featureVal.Value = featureId
		}

		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			limit = 10
		}
		offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil {
			offset = 0
		}

		fmt.Println("limit", limit, "offset", offset)

		data, err := redis.GetBannerGroup(*redisClient, tagVal, featureVal, limit, offset)
		if err == nil {
			render.JSON(w, r, data)
			return
		}

		banners, err := postgres.GetBannersFilteredByFeatureOrTagId(db, tagVal, featureVal, limit, offset)
		if err != nil {
			slog.Error("error while getting filtered banners")
			resp.Send500Error(w, r)
			return
		}
		redis.PutBannerGroup(*redisClient, tagVal, featureVal, banners, limit, offset)
		//w.WriteHeader(http.StatusOK)
		render.JSON(w, r, banners)
	}
}

func CreateBanner(redisClient *redis.Redis, db *postgres.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("token") != "admin_token" {
			resp.Send403Error(w, r)
			return
		}
		var createBanner models.CreateBannerRequest
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&createBanner)
		validate := validator.New()
		err2 := validate.Struct(createBanner)
		if err != nil || err2 != nil {
			slog.Error("error while parsing request body")
			resp.Send400Error(w, r)
			return
		}
		slog.Info("continue updating ", err, err2)
		banner, err := postgres.CreateUserBanner(db, createBanner)
		if err != nil {
			slog.Error("error while saving new banner")
			resp.Send500Error(w, r)
			return
		}

		slog.Info("success saved new banner " + string(banner.Id))
		for _, tagId := range createBanner.TagIds {
			redis.PutBanner(*redisClient, tagId, banner.FeatureId, banner)
		}
		//redis.PutBanner(*redisClient, )
		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, banner.Id)
	}
}

func UpdateBanner(redisClient *redis.Redis, db *postgres.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("token") != "admin_token" {
			resp.Send403Error(w, r)
			return
		}
		var bannerRequest models.CreateBannerRequest
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&bannerRequest)
		validate := validator.New()
		err2 := validate.Struct(bannerRequest)
		slog.Info("continue updating ", err, err2)
		if err != nil || err2 != nil {
			slog.Error("error while parsing request body")
			resp.Send400Error(w, r)
			return
		}
		slog.Info("continue updating ", err, err2)

		//fmt.Println(chi.URLParam(r, "id"))
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			slog.Error("error while parsing request id", err)
			resp.Send400Error(w, r)
			return
		}

		banner, err := postgres.GetBannerById(db, id)
		if err != nil {
			slog.Error("error while get banner by id ", id, err)
			resp.Send404Error(w, r)
			return
		}
		fmt.Println(banner)

		updatedBanner, err := postgres.UpdateUserBanner(db, id, bannerRequest, banner)
		if err != nil {
			slog.Error("error while updating request", err)
			resp.Send500Error(w, r)
			return
		}

		for _, tagId := range bannerRequest.TagIds {
			redis.PutBanner(*redisClient, tagId, banner.FeatureId, updatedBanner)
		}
		resp.Send200Success(w, r)
	}
}

func DeleteBanner(redisClient *redis.Redis, db *postgres.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("token") != "admin_token" {
			resp.Send403Error(w, r)
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			slog.Error("error while parsing request id", err)
			resp.Send400Error(w, r)
			return
		}

		_, err = postgres.GetBannerById(db, id)
		if err != nil {
			slog.Error("error while get banner by id ", id, err)
			resp.Send404Error(w, r)
			return
		}

		err = postgres.DeleteBannerById(db, id)
		if err != nil {
			slog.Error("error while deleting banner", id, err)
			resp.Send500Error(w, r)
			return
		}

		w.WriteHeader(http.StatusNoContent)

	}
}
