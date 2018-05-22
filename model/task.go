package model

import (
	"encoding/json"
	"fmt"
	"git-me/common"
	"git-me/utils"
	"io/ioutil"
	"sort"
)

type InnerTaskResp struct {
	Name  string `json:"name"`
	Pass  string `json:"pass"`
	Event int    `json:"event"`
	URL   string `json:"url"`
}

type InnerTaskReq struct {
	ErrCode int                 `json:"errCode"`
	Data    []*common.InnerTask `json:"data"`
}

func GetUserTask(p InnerTaskResp) (urls []string, err error) {
	body, err := json.Marshal(p)
	if err != nil {
		return
	}

	resp, err := utils.HttpPost(common.Host+common.ListRouter, body)
	if err != nil {
		return
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var req InnerTaskReq
	if err = json.Unmarshal(b, &req); err != nil {
		return
	}

	if req.ErrCode != 0 {
		err = fmt.Errorf("inner error, please contect developer")
		return
	}

	sort.Slice(req.Data, func(i, j int) bool {
		return req.Data[i].Sort > req.Data[j].Sort
	})

	for _, v := range req.Data {
		urls = append(urls, v.URL)
	}

	return
}

func UploadCurrentTaskStatus(p InnerTaskReq) (err error) {
	body, err := json.Marshal(p)
	if err != nil {
		return
	}

	resp, err := utils.HttpPost(common.Host+common.UploadRouter, body)
	if err != nil {
		return
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var req InnerTaskReq
	if err = json.Unmarshal(b, &req); err != nil {
		return
	}

	if req.ErrCode != 0 {
		err = fmt.Errorf("inner error, please contect developer")
		return
	}

	return
}
