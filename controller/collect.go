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
	"github.com/yusank/git-me/extractors"
	"github.com/yusank/git-me/models"

	"sort"

	"github.com/astaxie/beego/validation"
	"github.com/yusank/git-me/common"
)

type CollectController struct {
	BasicController
}

func (cc *CollectController) AddCollect() {
	uid := cc.GetString("userId")
	if uid == "" {
		cc.OnCustomError(consts.ErrNeedLogin)
		return
	}

	user, err := models.GetUserById(uid)
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

	if len(form.Formats) == 0 {
		cc.OnCustomError(consts.ErrNilToDownload)
		return
	}

	var (
		temp []common.FormatData
	)
	for _, v := range form.Formats {
		temp = append(temp, v)
	}

	sort.Slice(temp, func(i, j int) bool {
		return temp[i].Size > temp[j].Size
	})

	col = &models.CollectInfo{
		UserId:  user.Id,
		URL:     req.URL,
		Site:    form.Site,
		Size:    temp[0].Size,
		Title:   form.Title,
		Quality: temp[0].Quality,
	}

	if err = col.Insert(); err != nil {
		cc.OnError(err)
		return
	}

	cc.JSON("")
}

func (cc *CollectController) List() {
	uid := cc.GetString("userId")
	if uid == "" {
		cc.OnCustomError(consts.ErrNeedLogin)
		return
	}

	user, err := models.GetUserById(uid)
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
