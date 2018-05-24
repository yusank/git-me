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

	"git-me/consts"
	"git-me/models"

	"github.com/astaxie/beego/validation"
)

// 终端版本上交记录
type InnerController struct {
	BasicController
}

type InnerTaskReq struct {
	Name  string `json:"name" valid:"require"`
	Pass  string `json:"pass" valid:"require"`
	Event int    `json:"event" valid:"require"`
	URL   string `json:"url" valid:"require"`
}

// 更改任务的状态
func (ic *InnerController) HandleEvent() {
	var req InnerTaskReq
	if err := json.Unmarshal(ic.Ctx.Input.RequestBody, &req); err != nil {
		ic.OnError(err)
		return
	}

	valid := validation.Validation{}
	b, err := valid.Valid(&req)
	if err != nil {
		ic.OnError(err)
		return
	}

	if !b {
		ic.OnCustomError(consts.ErrInvalidParams)
		return
	}

	user := &models.User{
		Name: req.Name,
	}

	user, err = user.Get()
	if err != nil {
		ic.OnError(err)
		return
	}

	// 终端上传时，密码已加密过
	if user.Password != req.Pass {
		ic.OnCustomError(consts.ErrInvalidPass)
		return
	}

	task, err := models.GetTaskInfoByUserAndUrl(user.Id.Hex(), req.URL)
	if err != nil {
		ic.OnError(err)
		return
	}

	if task == nil {
		ic.OnCustomError(consts.ErrTaskNotFound)
		return
	}

	if task.Status == models.TaskStatusFinish {
		ic.OnCustomError(consts.ErrTaskFinish)
		return
	}

	task.Status = req.Event
	err = task.Update()

	if req.Event == models.TaskStatusFinish {
		his, err := models.GetHistory(user.Id.Hex(), req.URL)
		if err != nil {
			ic.OnError(err)
			return
		}

		if his == nil {
			his = &models.History{
				UserID: user.Id,
				URL:    req.URL,
			}

			err = his.Insert()
			goto finish
		}

		err = his.Update()
	}

finish:
	if err != nil {
		ic.OnError(err)
		return
	}

	ic.JSON("")
}

// 列出未完成任务
func (ic *InnerController) ListUserTasks() {
	var req InnerTaskReq
	err := json.Unmarshal(ic.Ctx.Input.RequestBody, &req)
	if err != nil {
		ic.OnError(err)
		return
	}

	user := &models.User{
		Name: req.Name,
	}

	user, err = user.Get()
	if err != nil {
		ic.OnError(err)
		return
	}

	// 上传时，用户信息已加密过
	if user.Password != req.Pass {
		ic.OnCustomError(consts.ErrInvalidPass)
		return
	}

	list, err := models.ListUnFinishedTaskInfo(user.Id.Hex())
	if err != nil {
		ic.OnError(err)
		return
	}

	ic.JSON(list)
}
