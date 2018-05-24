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
	"strconv"

	"github.com/yusank/git-me/consts"
	"github.com/yusank/git-me/models"

	"github.com/astaxie/beego/validation"
)

type TaskController struct {
	BasicController
}

type TaskInfo struct {
	Id     string `json:"id" valid:"Required"`
	UserId string `json:"userId"`
	URL    string `json:"url"`
	Sort   int    `json:"sort"`
}

func (tc *TaskController) ListTask() {
	uid := tc.GetSession(consts.SessionUserID)
	if uid == nil {
		tc.OnCustomError(consts.ErrNeedLogin)
		return
	}

	user, err := models.GetUserById(uid.(string))
	if err != nil {
		tc.OnError(err)
		return
	}

	if user == nil {
		tc.OnCustomError(consts.ErrUserNotFound)
		return
	}

	params := tc.Ctx.Input.Params()
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

	list, err := models.ListTaskInfo(user.Id.Hex(), defPage, defSize)
	if err != nil {
		tc.OnError(err)
		return
	}

	tc.JSON(list)
}

func (tc *TaskController) AddTask() {
	var req TaskInfo

	if err := json.Unmarshal(tc.Ctx.Input.RequestBody, &req); err != nil {
		tc.OnError(err)
		return
	}

	uid := tc.GetSession(consts.SessionUserID)
	if uid == nil {
		tc.OnCustomError(consts.ErrNeedLogin)
		return
	}

	user, err := models.GetUserById(uid.(string))
	if err != nil {
		tc.OnError(err)
		return
	}

	if user == nil {
		tc.OnCustomError(consts.ErrUserNotFound)
		return
	}

	task, err := models.GetTaskInfoByUserAndUrl(uid.(string), req.URL)
	if err != nil {
		tc.OnError(err)
		return
	}

	if task != nil {
		tc.OnCustomError(consts.ErrDataExists)
		return
	}

	task = &models.TaskInfo{
		UserId: user.Id,
		URL:    req.URL,
		Sort:   req.Sort,
	}

	if err := task.Insert(); err != nil {
		tc.OnError(err)
		return
	}

	tc.JSON("")
}

func (tc *TaskController) UpdateTask() {
	var req TaskInfo

	if err := json.Unmarshal(tc.Ctx.Input.RequestBody, &req); err != nil {
		tc.OnError(err)
		return
	}

	uid := tc.GetSession(consts.SessionUserID)
	if uid == nil {
		tc.OnCustomError(consts.ErrNeedLogin)
		return
	}

	user, err := models.GetUserById(uid.(string))
	if err != nil {
		tc.OnError(err)
		return
	}

	if user == nil {
		tc.OnCustomError(consts.ErrUserNotFound)
		return
	}

	v := validation.Validation{}
	b, err := v.Valid(&req)
	if err != nil {
		tc.OnError(err)
		return
	}

	if !b {
		tc.OnCustomError(consts.ErrInvalidParams)
		return
	}

	task, err := models.GetTaskInfoById(req.Id)
	if err != nil {
		tc.OnError(err)
		return
	}

	if task == nil {
		tc.OnCustomError(consts.ErrTaskNotFound)
	}

	if task.UserId != user.Id {
		tc.OnCustomError(consts.ErrNeedLogin)
		return
	}

	task.Sort = req.Sort

	if err = task.Update(); err != nil {
		tc.OnError(err)
		return
	}

	tc.JSON("")
}

func (tc *TaskController) DelTask() {
	var req TaskInfo

	if err := json.Unmarshal(tc.Ctx.Input.RequestBody, &req); err != nil {
		tc.OnError(err)
		return
	}

	uid := tc.GetSession(consts.SessionUserID)
	if uid == nil {
		tc.OnCustomError(consts.ErrNeedLogin)
		return
	}

	user, err := models.GetUserById(uid.(string))
	if err != nil {
		tc.OnError(err)
		return
	}

	if user == nil {
		tc.OnCustomError(consts.ErrUserNotFound)
		return
	}

	v := validation.Validation{}
	b, err := v.Valid(&req)
	if err != nil {
		tc.OnError(err)
		return
	}

	if !b {
		tc.OnCustomError(consts.ErrInvalidParams)
		return
	}

	task, err := models.GetTaskInfoById(req.Id)
	if err != nil {
		tc.OnError(err)
		return
	}

	if task == nil {
		tc.OnCustomError(consts.ErrTaskNotFound)
	}

	if task.UserId != user.Id {
		tc.OnCustomError(consts.ErrNeedLogin)
		return
	}

	if err = task.Delete(); err != nil {
		tc.OnError(err)
		return
	}

	tc.JSON("")
}
