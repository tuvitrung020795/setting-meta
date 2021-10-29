package handlers

import (
	"github.com/google/uuid"
	"github.com/praslar/lib/common"
	"github.com/tuvitrung020795/setting-meta/model"
	"github.com/tuvitrung020795/setting-meta/repo"
	"github.com/tuvitrung020795/setting-meta/utils"
	"gitlab.com/goxp/cloud0/ginext"
	"gitlab.com/goxp/cloud0/logger"
	"net/http"
)

type SettingMetaHandlers struct {
	obRepo *repo.SettingMeta
}

func NewSettingMetaHandlers(repo *repo.SettingMeta) *SettingMetaHandlers {
	return &SettingMetaHandlers{obRepo: repo}
}

func (h *SettingMetaHandlers) Create(r *ginext.Request) (*ginext.Response, error) {
	strUserID := ginext.GetUserID(r.GinCtx)
	owner, err := uuid.Parse(strUserID)
	if err != nil {
		return nil, ginext.NewError(http.StatusUnauthorized, err.Error())
	}

	req := model.SettingMetaRequest{}
	r.MustBind(&req)
	if err = common.CheckRequireValid(req); err != nil {
		return nil, ginext.NewError(http.StatusBadRequest, err.Error())
	}

	ob := &model.SettingMeta{
		BaseModel: model.BaseModel{
			CreatorID: &owner,
		},
	}
	common.Sync(req, ob)

	if err = h.obRepo.Create(r.Context(), ob); err != nil {
		return nil, ginext.NewError(http.StatusInternalServerError, err.Error())
	}

	return ginext.NewResponseData(http.StatusCreated, ob), nil
}

func (h *SettingMetaHandlers) List(r *ginext.Request) (*ginext.Response, error) {
	req := model.SettingMetaListRequest{}
	r.MustBind(&req)

	filter := &model.SettingMetaFilter{
		SettingMetaListRequest: req,
		Pager:                     ginext.NewPagerWithGinCtx(r.GinCtx),
	}

	result, err := h.obRepo.Filter(r.Context(), filter)
	if err != nil {
		return nil, ginext.NewError(http.StatusInternalServerError, err.Error())
	}
	r.MustNoError(err)

	resp := ginext.NewResponseWithPager(http.StatusOK, result.Records, result.Filter.Pager)
	return resp, nil
}

func (h *SettingMetaHandlers) GetOne(r *ginext.Request) (*ginext.Response, error) {
	ID := &uuid.UUID{}
	if ID = utils.ParseIDFromUri(r.GinCtx); ID == nil {
		return nil, ginext.NewError(http.StatusForbidden, "Wrong ID")
	}

	ob, err := h.obRepo.Get(r.Context(), *ID)
	if err != nil {
		return nil, ginext.NewError(http.StatusNotFound, err.Error())
	}
	r.MustNoError(err)

	return ginext.NewResponseData(http.StatusOK, ob), nil
}

func (h *SettingMetaHandlers) Update(r *ginext.Request) (*ginext.Response, error) {
	l := logger.WithCtx(r.Context(), "SettingMetaHandlers.Update")
	req := model.SettingMetaRequest{}
	strUserID := ginext.GetUserID(r.GinCtx)
	currentUser, err := uuid.Parse(strUserID)
	if err != nil {
		return nil, ginext.NewError(http.StatusUnauthorized, err.Error())
	}

	if req.ID = utils.ParseIDFromUri(r.GinCtx); req.ID == nil {
		return nil, ginext.NewError(http.StatusForbidden, "Wrong ID")
	}

	r.MustBind(&req)

	ob, err := h.obRepo.Get(r.Context(), *req.ID)
	if err != nil {
		l.WithError(err).WithField("todo.id", req.ID).Error("failed to query item")
		return nil, ginext.NewError(http.StatusNotFound, err.Error())
	}

	ob.UpdaterID = &currentUser
	common.Sync(req, ob)

	if err = h.obRepo.Update(r.Context(), ob); err != nil {
		return nil, ginext.NewError(http.StatusInternalServerError, err.Error())
	}

	return ginext.NewResponseData(http.StatusOK, ob), nil
}

func (h *SettingMetaHandlers) Delete(r *ginext.Request) (*ginext.Response, error) {
	strUserID := ginext.GetUserID(r.GinCtx)
	currentUser, err := uuid.Parse(strUserID)
	if err != nil {
		return nil, ginext.NewError(http.StatusUnauthorized, err.Error())
	}

	req := model.SettingMetaRequest{}
	if req.ID = utils.ParseIDFromUri(r.GinCtx); req.ID == nil {
		return nil, ginext.NewError(http.StatusForbidden, "Wrong ID")
	}

	ob, err := h.obRepo.Get(r.Context(), *req.ID)
	if err != nil {
		return nil, ginext.NewError(http.StatusForbidden, err.Error())
	}

	ob.UpdaterID = &currentUser
	r.MustNoError(h.obRepo.Delete(r.Context(), ob))

	return ginext.NewResponse(http.StatusOK), nil
}
