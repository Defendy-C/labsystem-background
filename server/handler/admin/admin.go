package admin

import (
	"github.com/gin-gonic/gin"
	"labsystem/model"
	"labsystem/model/srverr"
	"labsystem/server/handler"
	adminSrv "labsystem/service/admin"
	"net/http"
)

type HandlerAdmin struct {
	Srv adminSrv.ServiceAdmin
	handles []*handler.Handle
}

func (h *HandlerAdmin) RegisterAdminHandles(rg *gin.RouterGroup, authRg *gin.RouterGroup) {
	// the RouterGroup rg mustn't be Authenticate
	{
		rg.POST("/login", h.login)
	}
	// the RouterGroup hRg must be Authenticate
	{
		authRg.POST("/info", h.adminInfo)
		authRg.POST("/list", h.adminList)
		authRg.POST("/create", h.createAdmin)
	}
}

func (h *HandlerAdmin)login(ctx *gin.Context) {
	var req loginReq
	if err := ctx.Bind(&req); err != nil || !req.Valid() {
		ctx.JSON(http.StatusBadRequest, handler.NewResp(srverr.ErrInvalidParams, nil))
		return
	}
	t, err := h.Srv.Login(req.AdminNick, req.Password, req.Key, req.VCode)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, handler.NewResp(err, nil))
		return
	}

	ctx.JSON(http.StatusOK, handler.NewResp(nil, t))
}

func (h *HandlerAdmin)adminInfo(ctx *gin.Context) {
	uid, ok := ctx.Keys["uid"]
	if !ok {
		ctx.JSON(http.StatusForbidden, handler.NewResp(srverr.ErrForbidden, nil))
		return
	}
	if rid, ok := ctx.Keys["rid"]; !ok || int(rid.(float64)) != model.Administrator.Int() {
		ctx.JSON(http.StatusForbidden, handler.NewResp(srverr.ErrForbidden, nil))
		return
	}
	admin := h.Srv.QueryAdminById(uint(uid.(float64)))
	if admin == nil {
		ctx.JSON(http.StatusBadRequest, handler.NewResp(srverr.ErrSystemException ,nil))
		return
	}
	resp := new(InfoResp)
	resp.Name = admin.NickName
	resp.Powers = make([]*PowerOwner, len(model.PowerList))
	for i, v := range model.PowerList {
		resp.Powers[i] = &PowerOwner{
			Name: v.Name,
			Power: v.No,
			Own: admin.Power.Own(v.No),
		}
	}

	ctx.JSON(http.StatusOK, handler.NewResp(nil, resp))
}

func (h *HandlerAdmin)adminList(ctx *gin.Context) {
	var req *ListReq
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, handler.NewResp(srverr.ErrInvalidParams, nil))
		return
	}
	list, totalPage, totalCount := h.Srv.AdminList(&adminSrv.ListOpt{
		CreatedBy: req.CreatedBy,
	}, req.Page, req.PageSize)
	items := make([]*Item, len(list))
	for i, v := range list {
		items[i] = new(Item)
		items[i].ID = v.ID
		items[i].Name = v.NickName
		items[i].CreatedBy = v.CreatedBy
		items[i].CreatedAt = v.CreatedAt
		items[i].Power = make([]*PowerOwner, len(model.PowerList))
		for k, p := range model.PowerList {
			items[i].Power[k] = new(PowerOwner)
			items[i].Power[k].Power = p.No
			items[i].Power[k].Name = p.Name
			if v.Power.Own(p.No) {
				items[i].Power[k].Own = true
			}
		}
	}
	ctx.JSON(http.StatusOK, handler.NewResp(nil, &ListResp{
		List: items,
		TotalPage: totalPage,
		TotalCount: totalCount,
	}))
}

func (h *HandlerAdmin)createAdmin(ctx *gin.Context) {
	var req *CreateAdminReq
	if err := ctx.BindJSON(&req); err != nil || req.Valid() {
		ctx.JSON(http.StatusBadRequest, handler.NewResp(srverr.ErrInvalidParams, nil))
		return
	}

	p, _ := model.IntToPower(req.Power)
	if err := h.Srv.CreateAdmin(&model.Admin{
		NickName: req.Name,
		Password: req.Password,
		Power: p,
	}); err != nil {
		ctx.JSON(http.StatusBadRequest, handler.NewResp(srverr.ErrSystemException, nil))
		return
	}

	ctx.JSON(http.StatusOK, handler.NewResp(nil, nil))
}