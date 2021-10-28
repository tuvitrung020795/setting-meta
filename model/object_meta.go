package model

import (
	"github.com/google/uuid"
	"gitlab.com/goxp/cloud0/ginext"
)

type ObjectMeta struct {
	BaseModel
	SettingMetaId uuid.UUID `json:"setting_meta_id" gorm:"column:setting_meta_id"`
	ObjectId uuid.UUID `json:"object_id" gorm:"column:object_id"`
	ObjectType string `json:"object_type" gorm:"column:object_type;type:varchar(255)"`
	Value string `json:"value" gorm:"column:value"`
}

func (ObjectMeta) TableName() string {
	return "object_meta"
}

type ObjectMetaRequest struct {
	ID          *uuid.UUID `json:"id,omitempty"`
	SettingMetaId *uuid.UUID `json:"setting_meta_id" valid:"Required"`
	ObjectId *uuid.UUID `json:"object_id" valid:"Required"`
	ObjectType *string `json:"object_type"`
	Value *string `json:"value" valid:"Required"`
	
}

type ObjectMetaListRequest struct {
	CreatorID   *string `json:"creator_id" form:"creator_id"`
	SettingMetaId *uuid.UUID `json:"setting_meta_id" form:"setting_meta_id"`
	ObjectId *uuid.UUID `json:"object_id" form:"object_id"`
	ObjectType *string `json:"object_type" form:"object_type"`
}

type ObjectMetaFilter struct {
	ObjectMetaListRequest
	Pager *ginext.Pager
}

type ObjectMetaFilterResult struct {
	Filter  *ObjectMetaFilter
	Records []*ObjectMeta
}
