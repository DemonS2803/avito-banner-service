package mockdb

import (
	"avito-banner-service/internal/models"
	"avito-banner-service/internal/repositories/postgres"
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"time"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5433
	user     = "postgres"
	password = "postgres"
	dbname   = "avito_db"
)

func New() (*postgres.Storage, error) {
	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	fmt.Println("ready to connect psql")
	// open database
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	// check db
	err = db.Ping()
	CheckError(err)

	fmt.Println("Connected!")

	initQuery := readInitFile()
	rows, err := db.Query(initQuery)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	for rows.Next() {
		var version string
		rows.Scan(&version)
		fmt.Println(version)
	}
	fmt.Println("success read init sql file")

	fillQuery := readFillFile()
	rows, err = db.Query(fillQuery)
	if err != nil {
		fmt.Println("error while fill test data. continue with that data or rebuild postgres container")
		//log.Fatal(err)
		//panic(err)
	}

	return &postgres.Storage{Db: db}, nil
}

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}

func readInitFile() string {
	file, err := os.Open("resources/database/init.sql")
	CheckError(err)
	b, err := io.ReadAll(file)
	CheckError(err)
	//fmt.Println(b)
	return string(b)
}

func readFillFile() string {
	file, err := os.Open("resources/database/fill.sql")
	CheckError(err)
	b, err := io.ReadAll(file)
	CheckError(err)
	//fmt.Println(b)
	return string(b)
}

func GetBannersFilteredByFeatureOrTagId(storage *postgres.Storage, tagVal models.NilInt, featureVal models.NilInt) ([]models.UserBannerFilteredResponse, error) {
	rows, err := storage.Db.Query(`SELECT ub.id, ub.content, ub.is_active, ub.feature_id, ub.created_at, ub.updated_at, tag_id
				FROM user_banners_tags join user_banners ub on ub.id = user_banners_tags.banner_id
				where (feature_id = $1 OR $1 IS NULL) and (tag_id = $2 OR $2 IS NULL) ;`, featureVal.GetValue(), tagVal.GetValue())
	if err != nil {
		slog.Info("error while getting banners data", err)
		return []models.UserBannerFilteredResponse{}, err
	}

	//fmt.Println(rows)
	banners := []models.UserBannerFilteredResponse{}
	bannerIdPositionMap := make(map[int]int)
	for rows.Next() {
		//fmt.Println(rows)
		var banner models.UserBannerFilteredResponse
		var tagId int
		rows.Scan(&banner.BannerId, &banner.Content, &banner.IsActive, &banner.FeatureId, &banner.CreatedAt, &banner.UpdatedAt, &tagId)

		if bannerIdPositionMap[banner.BannerId] == 0 {
			banners = append(banners, banner)
			bannerIdPositionMap[banner.BannerId] = len(banners)
			banners[bannerIdPositionMap[banner.BannerId]-1].TagIds = []int{tagId}
		} else {
			banners[bannerIdPositionMap[banner.BannerId]-1].TagIds = append(banners[bannerIdPositionMap[banner.BannerId]-1].TagIds, tagId)
		}
	}
	//fmt.Println(bannerIdPositionMap)
	//fmt.Println(banners)
	return banners, nil
}

func CreateUserBanner(storage *postgres.Storage, createBanner models.CreateBannerRequest) (models.UserBanner, error) {

	var banner = models.UserBanner{
		Id:        GetNextUserBannerId(storage),
		Content:   createBanner.Content,
		IsActive:  createBanner.IsActive,
		FeatureId: createBanner.FeatureId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	slog.Info("new banner")
	fmt.Println(banner)

	ctx := context.Background()
	tx, err := storage.Db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO user_banners (id, content, is_active, feature_id, created_at, updated_at) values ($1, $2, $3, $4, $5, $6);`, banner.Id, banner.Content, banner.IsActive, banner.FeatureId, banner.CreatedAt, banner.UpdatedAt)
	if err != nil {
		slog.Error("error while saving new banner ", err)
		tx.Rollback()
		return models.UserBanner{}, err
	}

	for _, tagId := range createBanner.TagIds {
		_, err := tx.ExecContext(ctx, `INSERT INTO user_banners_tags (banner_id, tag_id) values ($1, $2);`, banner.Id, tagId)
		if err != nil {
			slog.Error("error adding tag id", tagId, err)
			tx.Rollback()
			return models.UserBanner{}, err
		}
	}

	tx.Commit()

	return banner, nil
}

func GetNextUserBannerId(storage *postgres.Storage) int {
	var maxId int
	err := storage.Db.QueryRow(`select MAX(id) FROM user_banners;`).Scan(&maxId)
	if err != nil {
		slog.Error("error while saving new banner ", err)
		return -1
	}
	return maxId + 1
}
