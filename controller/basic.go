package controller

import (
	"git-me/consts"
	"net/http"

	"github.com/astaxie/beego"
	"strings"
)

type BasicController struct {
	beego.Controller
}

func (this *BasicController) OnCustomError(e consts.ErrorType) {
	this.Ctx.ResponseWriter.WriteHeader(e.StatusCode)
	this.Data["json"] = map[string]interface{}{
		"errcode": e.ErrorCode,
		"errmsg":  e.ErrorMsg,
	}
	this.ServeJSON()
}

func (this *BasicController) JSON(data interface{}) {
	this.Ctx.ResponseWriter.WriteHeader(http.StatusOK)
	this.Data["json"] = map[string]interface{}{
		"errcode": 0,
		"data":    data,
	}
	this.ServeJSON()
}

func (this *BasicController) OnError(err error) {
	this.OnCustomError(consts.MakeError(err))
}

func (this *BasicController) GetIp() string {
	ips := this.Ctx.Input.Proxy()
	if len(ips) > 0 && ips[0] != "" {
		return ips[0]
	}

	ips = strings.Split(this.Ctx.Request.RemoteAddr, ":")
	if len(ips) > 0 {
		return ips[0]
	}

	return consts.LocalIp
}
