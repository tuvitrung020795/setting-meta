package model

import (
	"github.com/google/uuid"
	"gitlab.com/goxp/cloud0/ginext"
)

type SettingMeta struct {
	BaseModel
	Name string `json:"name" gorm:"column:name;type:varchar(255)"`
	Type string `json:"type" gorm:"column:type;type:varchar(255)"`
	Key string `json:"key" gorm:"column:key"`
	Required bool `json:"required" gorm:"column:required"`
	Description string `json:"description" gorm:"column:description"`
}

func (SettingMeta) TableName() string {
	return "setting_meta"
}

type SettingMetaRequest struct {
	ID          *uuid.UUID `json:"id,omitempty"`
	Name *string `json:"name" valid:"Required"`
	Type *string `json:"type" valid:"Required"`
	Key *string `json:"key" valid:"Required"`
	Required *bool `json:"required" valid:"Required"`
	Description *string `json:"description"`
}

type SettingMetaListRequest struct {
	CreatorID   *string `json:"creator_id" form:"creator_id"`
	Name *string `json:"name" form:"name"`
	Type *string `json:"type" form:"type"`
	Key *string `json:"key" form:"key"`
	Required *bool `json:"required" form:"required"`
}

type SettingMetaFilter struct {
	SettingMetaListRequest
	Pager *ginext.Pager
}

type SettingMetaFilterResult struct {
	Filter  *SettingMetaFilter
	Records []*SettingMeta
}
