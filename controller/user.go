/*
 * MIT License
 *
 * Copyright (c) 2018 Yusan Kurban <yusankurban@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

/*
 * Revision History:
 *     Initial: 2018/04/01        Yusan Kurban
 */

package controller

import (
	"encoding/json"

	"github.com/yusank/git-me/consts"
	"github.com/yusank/git-me/models"
	"github.com/yusank/git-me/utils"

	"fmt"

	"github.com/astaxie/beego/validation"
)

type UserController struct {
	BasicController
}

func (uc *UserController) Register() {
	var r models.UserRegister
	fmt.Println(string(uc.Ctx.Input.RequestBody))
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
		uc.OnCustomError(consts.ErrInvalidAccountOrPwd)
		return
	}

	uc.SetSession(consts.SessionUserID, user.Id.Hex())

	uc.Data["json"] = map[string]interface{}{"errcode": consts.ErrCodeSuccess, "data": user}
	uc.ServeJSON()
}

func (uc *UserController) GetInfo() {
	userID := uc.GetSession(consts.SessionUserID)
	if userID == nil {
		uc.OnCustomError(consts.ErrSessionNotFound)
		return
	}

	user, err := models.GetUserById(userID.(string))
	if err != nil {
		uc.OnError(err)
		return
	}

	uc.Data["json"] = map[string]interface{}{"errcode": consts.ErrCodeSuccess, "data": user}
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

	uc.JSON(oldInfo)
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
