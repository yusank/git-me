package controller

import (
	"encoding/json"

	"git-me/consts"
	"git-me/models"
	"git-me/utils"

	"github.com/astaxie/beego/validation"
)

type InnerTaskController struct {
	BasicController
}

type InnerTaskReq struct {
	Name  string `json:"name" valid:"require"`
	Pass  string `json:"pass" valid:"require"`
	Event int    `json:"event" valid:"require"`
	URL   string `json:"url" valid:"require"`
}

// 更改任务的状态
func (itc *InnerTaskController) HandleEvent() {
	var req InnerTaskReq
	if err := json.Unmarshal(itc.Ctx.Input.RequestBody, &req); err != nil {
		itc.OnError(err)
		return
	}

	valid := validation.Validation{}
	b, err := valid.Valid(&req)
	if err != nil {
		itc.OnError(err)
		return
	}

	if !b {
		itc.OnCustomError(consts.ErrInvalidParams)
		return
	}

	user := &models.User{
		Name: req.Name,
	}

	user, err = user.Get()
	if err != nil {
		itc.OnError(err)
		return
	}

	if user.Password != utils.StringMd5(req.Pass) {
		itc.OnCustomError(consts.ErrInvalidPass)
		return
	}

	task, err := models.GetTaskInfoByUserAndUrl(user.Id.Hex(), req.URL)
	if err != nil {
		itc.OnError(err)
		return
	}

	if task == nil {
		itc.OnCustomError(consts.ErrTaskNotFound)
		return
	}

	if task.Status == models.TaskStatusFinish {
		itc.OnCustomError(consts.ErrTaskFinish)
		return
	}

	task.Status = req.Event

	if err = task.Update(); err != nil {
		itc.OnError(err)
		return
	}

	itc.JSON("")
}

// 列出未完成任务
func (itc *InnerTaskController) ListUserTasks() {
	var req InnerTaskReq
	err := json.Unmarshal(itc.Ctx.Input.RequestBody, &req)
	if err != nil {
		itc.OnError(err)
		return
	}

	user := &models.User{
		Name: req.Name,
	}

	user, err = user.Get()
	if err != nil {
		itc.OnError(err)
		return
	}

	if user.Password != utils.StringMd5(req.Pass) {
		itc.OnCustomError(consts.ErrInvalidPass)
		return
	}

	list, err := models.ListUnFinishedTaskInfo(user.Id.Hex())
	if err != nil {
		itc.OnError(err)
		return
	}

	itc.JSON(list)
}
