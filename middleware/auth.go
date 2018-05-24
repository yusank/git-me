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

package middleware

import (
	"github.com/yusank/git-me/consts"
	"github.com/yusank/git-me/models"

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
	userID := ctx.Input.CruSession.Get(consts.SessionUserID)
	if userID == nil {
		authFailed(ctx)
		return
	}
	user, err := models.GetUserById(userID.(string))
	if err != nil || user == nil {
		authFailed(ctx)
		return
	}

	ctx.Input.SetData("user", user)
}