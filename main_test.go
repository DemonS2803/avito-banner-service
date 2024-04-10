package main

import (
	"avito-banner-service/internal/http-server/router"
	"avito-banner-service/internal/models"
	"avito-banner-service/internal/repositories/postgres"
	"avito-banner-service/internal/repositories/postgres/mockdb"
	"avito-banner-service/internal/repositories/redis/mock-cache"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPing(t *testing.T) {
	req, err := http.NewRequest("GET", "/ping", nil)
	req.Header.Set("token", "admin_token")
	if err != nil {
		t.Fatal(err)
	}

	newRecorder := httptest.NewRecorder()

	router.Routes(nil, nil).ServeHTTP(newRecorder, req)

	statusCode := 200
	if newRecorder.Result().StatusCode != statusCode {
		t.Errorf("TestPing() test returned an unexpected result: got %v want %v", newRecorder.Result().StatusCode, statusCode)
	}
}

func TestTokenAuth(t *testing.T) {
	req, err := http.NewRequest("GET", "/ping", nil)
	if err != nil {
		t.Fatal(err)
	}

	newRecorder := httptest.NewRecorder()

	router.Routes(nil, nil).ServeHTTP(newRecorder, req)

	statusCode := 401
	if newRecorder.Result().StatusCode != statusCode {
		t.Errorf("TestTokenAuth() test returned an unexpected result: got %v want %v", newRecorder.Result().StatusCode, statusCode)
	}
}

//func TestHandleGetUserBanner(t *testing.T) {
//
//	db, err := mockdb.New()
//	if err != nil {
//		slog.Error("failed to init storage")
//		os.Exit(1)
//	}
//
//	redisClient, err := redis.New()
//	if err != nil {
//		slog.Error("failed to init redis")
//		os.Exit(1)
//	}
//
//	req, err := http.NewRequest("GET", "/user_banner?tag_id=1&feature_id=1&use_last_revision=false", nil)
//	req.Header.Set("token", "user_token")
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	newRecorder := httptest.NewRecorder()
//
//	router.Routes(redisClient, db).ServeHTTP(newRecorder, req)
//
//	statusCode := 200
//	body, err := io.ReadAll(newRecorder.Body)
//	fmt.Println(body)
//	if newRecorder.Result().StatusCode != statusCode {
//		t.Errorf("HandleGetUserBannerTest() test returned an unexpected result: got %v want %v", newRecorder.Result().StatusCode, statusCode)
//	}
//}

func TestHandleGetUserBanner(t *testing.T) {

	db, err := mockdb.New()
	if err != nil {
		slog.Error("failed to init storage")
		t.Fatal(err)
	}

	var respContent string
	json.Unmarshal([]byte("{\n    \"url\": \"https://example.com/new-fashion-arrivals\",\n    \"text\": \"Discover the latest trends and styles for the season.\",\n    \"title\": \"New Arrivals in Fashion\"\n}"), &respContent)

	redisClient, err := mock_cache.New()
	if err != nil {
		slog.Error("failed to init redis")
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/user_banner?tag_id=1&feature_id=1&use_last_revision=false", nil)
	req.Header.Set("token", "user_token")
	if err != nil {
		t.Fatal(err)
	}

	newRecorder := httptest.NewRecorder()

	router.Routes(redisClient, db).ServeHTTP(newRecorder, req)

	statusCode := 200
	body, err := io.ReadAll(newRecorder.Body)

	var gotContent string
	json.Unmarshal(body, &gotContent)

	if newRecorder.Result().StatusCode != statusCode || gotContent != respContent {
		t.Errorf("HandleGetUserBannerTest() test returned an unexpected result: got %v want %v", newRecorder.Result().StatusCode, statusCode)
	}
}

func TestHandleGetBannerByAdmin(t *testing.T) {

	db, err := mockdb.New()
	if err != nil {
		slog.Error("failed to init storage")
		t.Fatal(err)
	}

	//banners, err := mockdb.GetBannersFilteredByFeatureOrTagId(db, models.NilInt{Null: true}, models.NilInt{Null: false, Value: 3})
	//if err != nil {
	//	t.Fatal("error when get test data from mockdb")
	//}

	//respContent, _ := json.Marshal(banners)

	redisClient, err := mock_cache.New()
	if err != nil {
		slog.Error("failed to init redis")
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/banner?feature_id=3", nil)
	req.Header.Set("token", "admin_token")
	if err != nil {
		t.Fatal(err)
	}

	newRecorder := httptest.NewRecorder()

	router.Routes(redisClient, db).ServeHTTP(newRecorder, req)

	statusCode := 200
	//body, err := io.ReadAll(newRecorder.Body)

	if newRecorder.Result().StatusCode != statusCode {
		t.Errorf("HandleGetUserBannerTest() test returned an unexpected result: got %v want %v", newRecorder.Result().StatusCode, statusCode)
	}
}

func TestCreateBanner(t *testing.T) {

	db, err := mockdb.New()
	if err != nil {
		slog.Error("failed to init storage")
		t.Fatal(err)
	}

	redisClient, err := mock_cache.New()
	if err != nil {
		slog.Error("failed to init redis")
		t.Fatal(err)
	}

	ids := []int{1, 4}
	featureId := 2
	content := `{"data": "some"}`
	isActive := true

	banner := models.CreateBannerRequest{
		TagIds:    ids,
		FeatureId: featureId,
		Content:   json.RawMessage(content),
		IsActive:  isActive,
	}

	body, err := json.Marshal(banner)

	req, err := http.NewRequest("POST", "/banner", strings.NewReader(string(body)))
	req.Header.Set("token", "admin_token")
	if err != nil {
		t.Fatal(err)
	}

	newRecorder := httptest.NewRecorder()

	router.Routes(redisClient, db).ServeHTTP(newRecorder, req)

	statusCode := 201

	if newRecorder.Result().StatusCode != statusCode {
		t.Errorf("HandleGetUserBannerTest() test returned an unexpected result: got %v want %v", newRecorder.Result().StatusCode, statusCode)
	}

	err = postgres.DeleteBannerById(db, 11)
	if err != nil {
		t.Fatal("created banner not deleted((")
	}
}

func TestPatchBanner(t *testing.T) {

	db, err := mockdb.New()
	if err != nil {
		slog.Error("failed to init storage")
		t.Fatal(err)
	}

	redisClient, err := mock_cache.New()
	if err != nil {
		slog.Error("failed to init redis")
		t.Fatal(err)
	}

	ids := []int{1, 4}
	featureId := 2
	content := `{"data": "some EDITED"}`
	isActive := true

	banner := models.CreateBannerRequest{
		TagIds:    ids,
		FeatureId: featureId,
		Content:   json.RawMessage(content),
		IsActive:  isActive,
	}

	body, err := json.Marshal(banner)

	req, err := http.NewRequest("PATCH", "/banner/4", strings.NewReader(string(body)))
	req.Header.Set("token", "admin_token")
	if err != nil {
		t.Fatal(err)
	}

	newRecorder := httptest.NewRecorder()

	router.Routes(redisClient, db).ServeHTTP(newRecorder, req)

	statusCode := 200

	if newRecorder.Result().StatusCode != statusCode {
		t.Errorf("HandleGetUserBannerTest() test returned an unexpected result: got %v want %v", newRecorder.Result().StatusCode, statusCode)
	}
}

func TestDeleteBanner(t *testing.T) {

	db, err := mockdb.New()
	if err != nil {
		slog.Error("failed to init storage")
		t.Fatal(err)
	}

	redisClient, err := mock_cache.New()
	if err != nil {
		slog.Error("failed to init redis")
		t.Fatal(err)
	}

	req, err := http.NewRequest("DELETE", "/banner/4", nil)
	req.Header.Set("token", "admin_token")
	if err != nil {
		t.Fatal(err)
	}

	newRecorder := httptest.NewRecorder()

	router.Routes(redisClient, db).ServeHTTP(newRecorder, req)

	statusCode := 204

	if newRecorder.Result().StatusCode != statusCode {
		t.Errorf("HandleGetUserBannerTest() test returned an unexpected result: got %v want %v", newRecorder.Result().StatusCode, statusCode)
	}

	ids := []int{1, 4}
	featureId := 2
	content := `{"data": "some"}`
	isActive := true

	banner := models.CreateBannerRequest{
		TagIds:    ids,
		FeatureId: featureId,
		Content:   json.RawMessage(content),
		IsActive:  isActive,
	}

	postgres.CreateUserBannerWithId(db, banner, 4)
	if err != nil {
		t.Fatal("create banner after delete error")
	}
}
