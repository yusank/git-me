package controller

import (
	"encoding/json"
	"git-me/consts"
	"git-me/models"
	"strconv"

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
