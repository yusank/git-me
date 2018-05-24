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
