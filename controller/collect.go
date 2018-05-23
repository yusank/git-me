package controller

import (
	"encoding/json"
	"git-me/consts"
	"git-me/models"
	"strconv"

	"git-me/extractors"

	"github.com/astaxie/beego/validation"
)

type CollectController struct {
	BasicController
}

func (cc *CollectController) AddCollect() {
	uid := cc.GetSession(consts.SessionUserID)
	if uid == nil {
		cc.OnCustomError(consts.ErrNeedLogin)
		return
	}

	user, err := models.GetUserById(uid.(string))
	if err != nil {
		cc.OnError(err)
		return
	}

	if user == nil {
		cc.OnCustomError(consts.ErrUserNotFound)
		return
	}

	var req DownloadInfo
	if err = json.Unmarshal(cc.Ctx.Input.RequestBody, &req); err != nil {
		cc.OnError(err)
		return
	}

	valid := validation.Validation{}
	ok, err := valid.Valid(&req)
	if err != nil {
		cc.OnError(err)
		return
	}

	if !ok {
		cc.OnCustomError(consts.ErrInvalidParams)
		return
	}

	col, err := models.GetCollectByUserID(user.Id.Hex(), req.URL)
	if err != nil {
		cc.OnError(err)
		return
	}

	if col != nil {
		cc.OnCustomError(consts.ErrDataExists)
		return
	}

	form, err := extractors.MatchUrl(req.URL)
	if err != nil || form == nil {
		cc.OnCustomError(consts.ErrNilToDownload)
		return
	}

	col = &models.CollectInfo{
		UserId: user.Id,
		URL:    req.URL,
		Site:   form.Site,
		Size:   form.Formats[0].Size,
	}

	if err = col.Insert(); err != nil {
		cc.OnError(err)
		return
	}

	cc.JSON("")
}

func (cc *CollectController) List() {
	uid := cc.GetSession(consts.SessionUserID)
	if uid == nil {
		cc.OnCustomError(consts.ErrNeedLogin)
		return
	}

	user, err := models.GetUserById(uid.(string))
	if err != nil {
		cc.OnError(err)
		return
	}

	if user == nil {
		cc.OnCustomError(consts.ErrUserNotFound)
		return
	}

	params := cc.Ctx.Input.Params()
	defPage := 1
	defSize := 10

	p, ok := params["page"]
	if ok {
		defPage, _ = strconv.Atoi(p)
	}

	s, ok := params["size"]
	if ok {
		defSize, _ = strconv.Atoi(s)
	}

	list, err := models.ListCollect(user.Id.Hex(), defPage, defSize)
	if err != nil {
		cc.OnError(err)
		return
	}

	cc.JSON(list)
}
