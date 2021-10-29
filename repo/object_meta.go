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

type ObjectMeta struct {
	db    *gorm.DB
	debug bool
}

func (r *ObjectMeta) DBWithTimeout(ctx context.Context) (*gorm.DB, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(ctx, generalQueryTimeout)
	return r.db.WithContext(ctx), cancel
}

func NewObjectMetaRepo(db *gorm.DB) *ObjectMeta {
	return &ObjectMeta{db: db}
}

func (r *ObjectMeta) Create(ctx context.Context, ob *model.ObjectMeta) error {
	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	if utils.MustExist(tx, &model.SettingMeta{}, "id", ob.SettingMetaId) {
		return fmt.Errorf("Invalid Setting meta")
	}

	return tx.Save(ob).Error
}

func (r *ObjectMeta) Get(ctx context.Context, id uuid.UUID) (*model.ObjectMeta, error) {
	o := &model.ObjectMeta{}

	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	err := tx.First(&o, id).Error
	return o, err
}

func (r *ObjectMeta) Filter(ctx context.Context, f *model.ObjectMetaFilter) (*model.ObjectMetaFilterResult, error) {
	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	log := logger.WithCtx(ctx, "ObjectMeta.Filter")
	tx = tx.WithContext(ctx).Model(&model.ObjectMeta{})

	op := tx.Where
	tx = utils.FilterIfNotNil(f.CreatorID, tx, op, "creator_id = ?")
	tx = utils.FilterIfNotNil(f.SettingMetaId, tx, op, "setting_meta_id = ?")
	tx = utils.FilterIfNotNil(f.ObjectId, tx, op, "object_id = ?")
	tx = utils.FilterIfNotNil(f.ObjectType, tx, op, "object_type = ?")
	result := &model.ObjectMetaFilterResult{
		Filter:  f,
		Records: []*model.ObjectMeta{},
	}

	f.Pager.SortableFields = []string{"id", "created_at", "updated_at"}
	pager := result.Filter.Pager

	tx = pager.DoQuery(&result.Records, tx)
	if tx.Error != nil {
		log.WithError(tx.Error).Error("error while filter object meta")
	}

	return result, tx.Error
}

func (r *ObjectMeta) Update(ctx context.Context, update *model.ObjectMeta) error {
	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	if utils.MustExist(tx, &model.SettingMeta{}, "id", update.SettingMetaId) {
		return fmt.Errorf("Invalid Setting meta")
	}
	return tx.WithContext(ctx).Where("id = ?", update.ID).Save(&update).Error
}

func (r *ObjectMeta) Delete(ctx context.Context, d *model.ObjectMeta) error {
	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	return tx.WithContext(ctx).Delete(&d).Error
}

func (r *ObjectMeta) GetOneFlexible(ctx context.Context, field string, value interface{}) (model.ObjectMeta, error) {
	o := model.ObjectMeta{}

	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	err := tx.Where(field+" = ? ", value).First(&o).Error
	return o, err
}

func (r *ObjectMeta) GetValue(ctx context.Context,settingMetaId uuid.UUID,objectId uuid.UUID) (*model.ObjectMeta, error)  {
	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	o := &model.ObjectMeta{}

	if err := tx.Model(&model.ObjectMeta{}).Where("setting_meta_id = ? AND object_id = ?", settingMetaId,objectId).First(&o).Error; err != nil {
		return nil,err
	}

	return o,nil
}