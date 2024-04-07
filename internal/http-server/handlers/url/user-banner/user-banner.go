package user_banner

import (
	"avito-banner-service/internal/configs/postgres"
	"avito-banner-service/internal/configs/redis"
	"avito-banner-service/internal/models"
	resp "avito-banner-service/internal/utils/response"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
)

type UserBannerRequest struct {
	TagId          int    `json:"tag_id" validate:"required"`
	FeatureId      int    `json:"feature_id" validate:"required"`
	UseLastVersion bool   `json:"use_last_version" default:"false"`
	Token          string `json:"token" default:"user_token"`
}

type Response struct {
	resp.Response
	banner models.UserBanner
}

func GetBannerById(redisClient *redis.Redis, db *postgres.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

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

		//fmt.Println("useLastRevision ")
		//fmt.Println(useLastRevision)

		var banner models.UserBanner
		isCached := redis.GetBannerById(*redisClient, tagId, featureId, &banner)
		if useLastRevision || !isCached {
			banner, err = postgres.GetUserBannerByTagIdAndFeatureId(db, tagId, featureId)
			if err != nil {
				slog.Error("error while get banner from db: ", err)
				resp.Send404Error(w, r)
				return
			}
		}

		redis.PutBanner(*redisClient, tagId, featureId, banner)

		if !banner.IsActive && r.Header.Get("token") != "admin_token" {
			resp.Send403Error(w, r)
			return
		}

		//slog.Log(context.Background(), 0, fmt.Sprintf("id %d in system: %s", tagId, featureId, banner.Content))

		//w.Header().Set("Content-Type", "application/json")

		//res, err := json.Marshal(banner.Content)
		//fmt.Println(string(res))
		//fmt.Println("next id ", postgres.GetNextUserBannerId(db))
		render.JSON(w, r, banner.Content)
	}
}
