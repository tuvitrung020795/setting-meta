package repo

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/tuvitrung020795/setting-meta/model"
	"github.com/tuvitrung020795/setting-meta/utils"
	"gitlab.com/goxp/cloud0/logger"
	"gorm.io/gorm"
)

type SettingMetaDetail struct {
	db    *gorm.DB
	debug bool
}

func (r *SettingMetaDetail) DBWithTimeout(ctx context.Context) (*gorm.DB, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(ctx, generalQueryTimeout)
	return r.db.WithContext(ctx), cancel
}

func NewSettingMetaDetailRepo(db *gorm.DB) *SettingMetaDetail {
	return &SettingMetaDetail{db: db}
}

func (r *SettingMetaDetail) Create(ctx context.Context, ob *model.SettingMetaDetail) error {
	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	if utils.MustExist(tx, &model.SettingMeta{}, "id", ob.SettingMetaId) {
		return fmt.Errorf("Invalid Setting Meta")
	}

	if r.CheckExistKey(tx, uuid.Nil,ob.SettingMetaId, ob.Key) {
		return fmt.Errorf("Detail key already exist")
	}

	return tx.Create(ob).Error
}

func (r *SettingMetaDetail) Get(ctx context.Context, id uuid.UUID) (*model.SettingMetaDetail, error) {
	o := &model.SettingMetaDetail{}

	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	err := tx.First(&o, id).Error
	return o, err
}

func (r *SettingMetaDetail) Filter(ctx context.Context, f *model.SettingMetaDetailFilter) (*model.SettingMetaDetailFilterResult, error) {
	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	log := logger.WithCtx(ctx, "SettingMetaDetail.Filter")
	tx = tx.WithContext(ctx).Model(&model.SettingMetaDetail{})

	op := tx.Where
	tx = utils.FilterIfNotNil(f.CreatorID, tx, op, "creator_id = ?")
	tx = utils.FilterIfNotNil(f.Key, tx, op, "key = ?")
	tx = utils.FilterIfNotNil(f.SettingMetaId, tx, op, "setting_meta_id = ?")
	result := &model.SettingMetaDetailFilterResult{
		Filter:  f,
		Records: []*model.SettingMetaDetail{},
	}

	f.Pager.SortableFields = []string{"id", "created_at", "updated_at"}
	pager := result.Filter.Pager

	tx = pager.DoQuery(&result.Records, tx)
	if tx.Error != nil {
		log.WithError(tx.Error).Error("error while filter object meta")
	}

	return result, tx.Error
}

func (r *SettingMetaDetail) Update(ctx context.Context, ob *model.SettingMetaDetail) error {
	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	if utils.MustExist(tx, &model.SettingMeta{}, "id", ob.SettingMetaId) {
		return fmt.Errorf("Invalid Setting Meta")
	}

	if r.CheckExistKey(tx, ob.ID,ob.SettingMetaId, ob.Key) {
		return fmt.Errorf("Detail key already exist")
	}

	return tx.WithContext(ctx).Where("id = ?", ob.ID).Save(&ob).Error
}

func (r *SettingMetaDetail) Delete(ctx context.Context, d *model.SettingMetaDetail) error {
	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	return tx.WithContext(ctx).Delete(&d).Error
}

func (r *SettingMetaDetail) GetOneFlexible(ctx context.Context, field string, value interface{}) (model.SettingMetaDetail, error) {
	o := model.SettingMetaDetail{}

	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	err := tx.Where(field+" = ? ", value).First(&o).Error
	return o, err
}

func (r *SettingMetaDetail) CheckExistKey(tx *gorm.DB, id uuid.UUID, settingMetaId uuid.UUID,key string) bool {
	var count int64 = 0

	if id != uuid.Nil {
		tx = tx.Where("id != ? ",id)
	}

	tx = tx.Where("setting_meta_id = ? AND key = ?", settingMetaId,key)
	err := tx.Model(&model.SettingMetaDetail{}).Count(&count).Error
	if err != nil {
		return true
	}
	return count > 0
}