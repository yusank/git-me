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

package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"

	"github.com/yusank/git-me/common"
	"github.com/yusank/git-me/utils"
)

type InnerTaskResp struct {
	Name     string  `json:"name"`
	Pass     string  `json:"pass"`
	Event    int     `json:"event"`
	Schedule float64 `json:"schedule"`
	URL      string  `json:"url"`
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
	fmt.Println(string(b))
	if err = json.Unmarshal(b, &req); err != nil {
		return
	}

	if req.ErrCode != 0 {
		fmt.Println(req)
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

func UploadCurrentTaskStatus(p InnerTaskResp) (err error) {
	if p.Name == "" || p.Pass == "" {
		return
	}

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
		fmt.Println(req)
		err = fmt.Errorf("inner error, please contect developer")
		return
	}

	return
}
