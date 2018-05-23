package controller

import (
	"git-me/consts"
	"git-me/models"
	"strconv"
)

type CollectController struct {
	BasicController
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
