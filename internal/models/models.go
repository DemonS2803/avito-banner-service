package models

import (
	"encoding/json"
	"time"
)

type Feature struct {
	Id   int    `json:"id" redis:"id"`
	Name string `json:"name" redis:"name"`
}

type UserBanner struct {
	Id        int             `json:"id" redis:"id"`
	Content   json.RawMessage `json:"content,omitempty" redis:"content"`
	IsActive  bool            `json:"is_active" redis:"is_active"`
	FeatureId int             `json:"feature_id,omitempty" redis:"feature_id"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type CreateBannerRequest struct {
	TagIds    []int           `json:"tag_ids" validate:"required"`
	FeatureId int             `json:"feature_id" validate:"required"`
	Content   json.RawMessage `json:"content" validate:"required"`
	IsActive  bool            `json:"is_active"`
}

type UserBannerFilteredResponse struct {
	BannerId  int             `json:"banner_id"`
	TagIds    []int           `json:"tag_ids"`
	FeatureId int             `json:"feature_id,omitempty"`
	Content   json.RawMessage `json:"content,omitempty"`
	IsActive  bool            `json:"is_active"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type NilInt struct {
	Value int
	Null  bool
}

func (n *NilInt) GetValue() interface{} {
	if n.Null {
		return nil
	}
	return n.Value
}
