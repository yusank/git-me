package controller

import (
	"git-me/consts"
	"git-me/models"
	"strconv"
)

type HistoryController struct {
	BasicController
}

func (hc *HistoryController) List() {
	uid := hc.GetSession(consts.SessionUserID)
	if uid == nil {
		hc.OnCustomError(consts.ErrNeedLogin)
		return
	}

	user, err := models.GetUserById(uid.(string))
	if err != nil {
		hc.OnError(err)
		return
	}

	if user == nil {
		hc.OnCustomError(consts.ErrUserNotFound)
		return
	}

	params := hc.Ctx.Input.Params()
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

	list, err := models.ListHistory(user.Id.Hex(), defPage, defSize)
	if err != nil {
		hc.OnError(err)
		return
	}

	hc.JSON(list)
}
