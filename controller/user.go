package controller

import (
	"encoding/json"
	"git-me/consts"
	"git-me/models"
	"git-me/utils"

	"github.com/astaxie/beego/validation"
)

type UserController struct {
	BasicController
}

func (uc *UserController) Register() {
	var r models.UserRegister
	err := json.Unmarshal(uc.Ctx.Input.RequestBody, &r)
	if err != nil {
		uc.OnError(err)
		return
	}

	if !utils.IsValidAccount(r.Name) {
		uc.OnCustomError(consts.ErrInvalidUserName)
		return
	}

	if !utils.IsValidEmail(r.Email) {
		uc.OnCustomError(consts.ErrInvalidEmail)
		return
	}

	if !utils.IsValidPassword(r.Pass) {
		uc.OnCustomError(consts.ErrInvalidPass)
		return
	}

	user := models.User{
		Name:     r.Name,
		Email:    r.Email,
		Password: utils.StringMd5(r.Pass),
	}

	if err = user.Insert(); err != nil {
		uc.OnError(err)
		return
	}

	uc.JSON("success")
}

func (uc *UserController) Login() {
	var l models.UserLogin

	err := json.Unmarshal(uc.Ctx.Input.RequestBody, &l)
	if err != nil {
		uc.OnError(err)
		return
	}

	user := &models.User{
		Name: l.Name,
	}

	user, err = user.Get()
	if err != nil {
		uc.OnError(err)
	}

	passMd5 := utils.StringMd5(l.Pass)
	if passMd5 != user.Password {
		uc.OnCustomError(consts.ErrInvalidAcountOrPwd)
		return
	}

	uc.SetSession(consts.SessionUserID, user.Id)

	uc.Data["json"] = map[string]interface{}{"code": consts.ErrCodeSuccess, "data": user}
	uc.ServeJSON()
}

func (uc *UserController) UpdateInfo() {
	var u models.User

	err := json.Unmarshal(uc.Ctx.Input.RequestBody, &u)
	if err != nil {
		uc.OnError(err)
		return
	}

	if u.Nickname == "" && u.HeadImg == "" {
		uc.JSON("success")
		return
	}

	oldInfo := &models.User{
		Name: u.Name,
	}

	oldInfo, err = oldInfo.Get()
	if err != nil {
		uc.OnError(err)
		return
	}

	if oldInfo == nil {
		uc.OnCustomError(consts.ErrUserNotFound)
	}

	if u.Nickname != "" {
		oldInfo.Nickname = u.Nickname
	}
	if u.HeadImg != "" {
		oldInfo.HeadImg = u.HeadImg
	}

	if err = oldInfo.Update(); err != nil {
		uc.OnError(err)
		return
	}

	uc.JSON("success")
}

func (uc *UserController) UpdatePass() {
	var up models.UpdatePass

	err := json.Unmarshal(uc.Ctx.Input.RequestBody, &up)
	if err != nil {
		uc.OnError(err)
		return
	}

	valid := validation.Validation{}
	b, err := valid.Valid(&up)
	if err != nil {
		uc.OnError(err)
		return
	}

	if !b {
		uc.OnCustomError(consts.ErrInvalidParams)
		return
	}

	oldInfo := &models.User{
		Name: up.Name,
	}

	oldInfo, err = oldInfo.Get()
	if err != nil {
		uc.OnError(err)
		return
	}

	if oldInfo == nil {
		uc.OnCustomError(consts.ErrUserNotFound)
		return
	}

	if utils.StringMd5(up.OldPass) != oldInfo.Password {
		uc.OnCustomError(consts.ErrInvalidPass)
		return
	}

	oldInfo.Password = utils.StringMd5(up.NewPass)

	if err := oldInfo.Update(); err != nil {
		uc.OnError(err)
		return
	}

	uc.JSON("success")
}

func (uc *UserController) Logout() {
	userID := uc.GetSession(consts.SessionUserID)
	if userID == nil {
		uc.OnCustomError(consts.ErrSessionNotFound)
		return
	}

	uc.DestroySession()
	uc.JSON("success")
}
