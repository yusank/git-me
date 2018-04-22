package middleware

import (
	"git-me/consts"

	"github.com/astaxie/beego/context"
)

func authFailed(ctx *context.Context) {
	e := consts.ErrNeedLogin
	ctx.Output.SetStatus(e.StatusCode)
	ctx.Output.JSON(map[string]interface{}{
		"errcode": e.ErrorCode,
		"errmsg":  e.ErrorMsg,
	}, true, false)
}

func AuthLogin(ctx *context.Context) {
	channelName := ctx.Input.CruSession.Get(consts.SessionUserID)
	if channelName == nil {
		authFailed(ctx)
		return
	}

	//ctx.Input.SetData("vendor", vendor)
}
