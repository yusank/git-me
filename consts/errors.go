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
	ErrInvalidParams      = NewError(200, 400, "参数错误")
	ErrUserNotFound       = NewError(200, 400, "用户未找到")
	ErrInvalidUserName    = NewError(200, 400, "非法的用户名")
	ErrInvalidEmail       = NewError(200, 400, "非法的邮箱")
	ErrInvalidPass        = NewError(200, 400, "非法的密码")
	ErrSessionNotFound    = NewError(200, 403, "session不存在")
	ErrInvalidAcountOrPwd = NewError(200, 400, "账号或密码错误")
	ErrBsonObjectId       = NewError(200, 400, "非法id")
	ErrSign               = NewError(200, 400, "签名错误")
	ErrNilToDownload      = NewError(200, 400, "没有可下载的信息")
	ErrNotFound           = NewError(200, 400, "搜索结果为空")

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
