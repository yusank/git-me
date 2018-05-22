package controller

import (
	"encoding/json"
	"git-me/consts"
	"git-me/models"
)

type TaskController struct {
	BasicController
}

type TaskInfo struct {
	UserId string `json:"userId"`
	URL    string `json:"url"`
	Sort   int    `json:"sort"`
}

func (tc *TaskController) ListTask() {
	var req TaskInfo

	if err := json.Unmarshal(tc.Ctx.Input.RequestBody, &req); err != nil {
		tc.OnError(err)
		return
	}

	list, err := models.ListTaskInfo(req.UserId)
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
		tc.OnCustomError(consts.ErrTaskExists)
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

	task, err := models.GetTaskInfoByUserAndUrl(uid.(string), req.URL)
	if err != nil {
		tc.OnError(err)
		return
	}

	if task == nil {
		tc.OnCustomError(consts.ErrTaskNotFound)
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

	task, err := models.GetTaskInfoByUserAndUrl(uid.(string), req.URL)
	if err != nil {
		tc.OnError(err)
		return
	}

	if task == nil {
		tc.OnCustomError(consts.ErrTaskNotFound)
	}

	if err = task.Delete(); err != nil {
		tc.OnError(err)
		return
	}

	tc.JSON("")
}
