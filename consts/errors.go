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

package consts

type ErrorType struct {
	StatusCode int
	ErrorCode  int
	ErrorMsg   string
}

func NewError(status, errcode int, errmsg string) (err ErrorType) {
	return ErrorType{
		StatusCode: status,
		ErrorCode:  errcode,
		ErrorMsg:   errmsg,
	}
}

//不自带code 是因为别的code都自定义  只有服务器错误不自定义
func MakeError(err error) (e ErrorType) {
	return ErrorType{
		StatusCode: 200,
		ErrorCode:  500,
		ErrorMsg:   err.Error(),
	}
}

var (
	ErrInvalidParams       = NewError(200, 400, "参数错误")
	ErrInvalidUrl          = NewError(200, 400, "非法的 url")
	ErrUserNotFound        = NewError(200, 400, "用户未找到")
	ErrInvalidUserName     = NewError(200, 400, "非法的用户名")
	ErrInvalidEmail        = NewError(200, 400, "非法的邮箱")
	ErrInvalidPass         = NewError(200, 400, "非法的密码")
	ErrSessionNotFound     = NewError(200, 403, "用户 id 为空")
	ErrInvalidAccountOrPwd = NewError(200, 400, "账号或密码错误")
	ErrIpCannotFound       = NewError(200, 400, "ip 无法获取")
	ErrSign                = NewError(200, 400, "签名错误")
	ErrNilToDownload       = NewError(200, 400, "没有可下载的信息")
	ErrNotFound            = NewError(200, 400, "搜索结果为空")
	ErrDataExists          = NewError(200, 400, "重复添加")
	ErrTaskNotFound        = NewError(200, 400, "任务无法找打")
	ErrTaskFinish          = NewError(200, 400, "任务已完成")
	ErrIpHasBlocked        = NewError(200, 0400, "该 ip 已限制使用")

	ErrNeedLogin = NewError(200, 403, "请先登录")
)

func MakeResponse(data interface{}) (ret *map[string]interface{}) {
	ret = &map[string]interface{}{
		"errcode": 0,
		"errmsg":  "",
		"data":    data,
	}
	return
}
