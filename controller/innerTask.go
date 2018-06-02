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

	"github.com/yusank/git-me/common"
	"github.com/yusank/git-me/consts"
	"github.com/yusank/git-me/models"

	"sort"

	"github.com/yusank/git-me/extractors"
)

// 终端版本上交记录
type InnerController struct {
	BasicController
}

type InnerTaskReq struct {
	Name     string  `json:"name" valid:"Required"`
	Pass     string  `json:"pass" valid:"Required"`
	Event    int     `json:"event" valid:"Required"`
	Schedule float64 `json:"schedule"`
	URL      string  `json:"url" valid:"Required"`
}

// 更改任务的状态
func (ic *InnerController) HandleEvent() {
	var (
		req InnerTaskReq
		err error
	)

	if err := json.Unmarshal(ic.Ctx.Input.RequestBody, &req); err != nil {
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

	if req.Event == models.TaskStatusDownloading {
		task.Schedule = req.Schedule
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
			vid, err := extractors.MatchUrl(req.URL)
			if err != nil {
				ic.OnError(err)
				return
			}

			if vid == nil {
				ic.OnCustomError(consts.ErrNilToDownload)
				return
			}

			var (
				temp []common.FormatData
			)
			for _, v := range vid.Formats {
				temp = append(temp, v)
			}

			sort.Slice(temp, func(i, j int) bool {
				return temp[i].Size > temp[j].Size
			})

			info := temp[0]
			his = &models.History{
				UserID:  user.Id,
				URL:     req.URL,
				Size:    info.Size,
				Site:    vid.Site,
				Type:    vid.Type,
				Title:   vid.Title,
				Quality: info.Quality,
			}

			err = his.Insert()
		} else {
			err = his.Update()
		}

		models.AddStatSite(his.Site)
	}

	if err != nil {
		ic.OnError(err)
		return
	}

	ic.JSON(nil)
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
