package controller

import (
	"encoding/json"

	"git-me/consts"
	"git-me/models"
	"git-me/utils"

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

	if user.Password != utils.StringMd5(req.Pass) {
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

	if user.Password != utils.StringMd5(req.Pass) {
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
