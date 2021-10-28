package model

import (
	"github.com/google/uuid"
	"gitlab.com/goxp/cloud0/ginext"
)

type SettingMetaDetail struct {
	BaseModel
	SettingMetaId uuid.UUID `json:"setting_meta_id"`
	Key string `json:"key"`
	Value string `json:"value"`
}

func (SettingMetaDetail) TableName() string {
	return "setting_meta_detail"
}

type SettingMetaDetailRequest struct {
	ID          *uuid.UUID `json:"id,omitempty"`
	SettingMetaId *uuid.UUID `json:"setting_meta_id" valid:"Required"`
	Key *string `json:"key" valid:"Required"`
	Value *string `json:"value"`
}

type SettingMetaDetailListRequest struct {
	CreatorID   *string `json:"creator_id" form:"creator_id"`
	SettingMetaId *uuid.UUID `json:"setting_meta_id" form:"setting_meta_id"`
	Key *string `json:"key" form:"key"`
	Value *string `json:"value" form:"value"`
}

type SettingMetaDetailFilter struct {
	SettingMetaDetailListRequest
	Pager *ginext.Pager
}

type SettingMetaDetailFilterResult struct {
	Filter  *SettingMetaDetailFilter
	Records []*SettingMetaDetail
}
