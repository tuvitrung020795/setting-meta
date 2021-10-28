package repo

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"gitlab.com/goxp/cloud0/logger"
	"gorm.io/gorm"
	"setting-meta/model"
	"setting-meta/utils"
)

type SettingMeta struct {
	db    *gorm.DB
	debug bool
}

func (r *SettingMeta) DBWithTimeout(ctx context.Context) (*gorm.DB, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(ctx, generalQueryTimeout)
	return r.db.WithContext(ctx), cancel
}

func NewSettingMetaRepo(db *gorm.DB) *SettingMeta {
	return &SettingMeta{db: db}
}

func (r *SettingMeta) Create(ctx context.Context, ob *model.SettingMeta) error {
	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	if  r.CheckExistKey(tx, uuid.Nil, ob.Key) {
		return fmt.Errorf("Key already exist")
	}

	return tx.Create(ob).Error
}

func (r *SettingMeta) Get(ctx context.Context, id uuid.UUID) (*model.SettingMeta, error) {
	o := &model.SettingMeta{}

	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	err := tx.First(&o, id).Error
	return o, err
}

func (r *SettingMeta) Filter(ctx context.Context, f *model.SettingMetaFilter) (*model.SettingMetaFilterResult, error) {
	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	log := logger.WithCtx(ctx, "SettingMeta.Filter")
	tx = tx.WithContext(ctx).Model(&model.SettingMeta{})

	op := tx.Where
	tx = utils.FilterIfNotNil(f.CreatorID, tx, op, "creator_id = ?")
	tx = utils.FilterIfNotNil(f.Key, tx, op, "key = ?")
	tx = utils.FilterIfNotNil(f.Name, tx, op, "name ILike ?")
	tx = utils.FilterIfNotNil(f.Type, tx, op, "type = ?")
	tx = utils.FilterIfNotNil(f.Required, tx, op, "required = ?")
	result := &model.SettingMetaFilterResult{
		Filter:  f,
		Records: []*model.SettingMeta{},
	}

	f.Pager.SortableFields = []string{"id", "created_at", "updated_at"}
	pager := result.Filter.Pager

	tx = pager.DoQuery(&result.Records, tx)
	if tx.Error != nil {
		log.WithError(tx.Error).Error("error while filter object meta")
	}

	return result, tx.Error
}

func (r *SettingMeta) Update(ctx context.Context, update *model.SettingMeta) error {
	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	if  r.CheckExistKey(tx, update.ID, update.Key) {
		return fmt.Errorf("Key already exist")
	}

	return tx.WithContext(ctx).Where("id = ?", update.ID).Save(&update).Error
}

func (r *SettingMeta) Delete(ctx context.Context, d *model.SettingMeta) error {
	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	return tx.WithContext(ctx).Delete(&d).Error
}

func (r *SettingMeta) GetOneFlexible(ctx context.Context, field string, value interface{}) (model.SettingMeta, error) {
	o := model.SettingMeta{}

	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	err := tx.Where(field+" = ? ", value).First(&o).Error
	return o, err
}

func (r *SettingMeta) CheckExistKey(tx *gorm.DB, id uuid.UUID,key string) bool {
	var count int64 = 0

	if id != uuid.Nil {
		tx = tx.Where("id != ? ",id)
	}

	err := tx.Model(&model.SettingMeta{}).Where("key = ?", key).Count(&count).Error
	if err != nil {
		return true
	}
	return count > 0
}