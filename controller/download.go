package controller

import (
	"fmt"
	"git-me/consts"
	"git-me/extractors"
	"git-me/models"
	"git-me/utils"
	"io"
	"os"
)

type DownloaderController struct {
	BasicController
}

type DownloadInfo struct {
	URL string `json:"url"`
}

func (dc *DownloaderController) ParseVideo() {
	var di DownloadInfo
	//if err := json.Unmarshal(dc.Ctx.Input.RequestBody, &di); err != nil {
	//	dc.OnError(err)
	//	return
	//}

	di.URL = "https://www.bilibili.com/video/av23606332/?spm_id_from=333.334.bili_game.4"
	uid := dc.GetSession(consts.SessionUserID)
	if uid == nil {
		ip := dc.GetIp()
		if ip == "" {
			dc.OnCustomError(consts.ErrIpCannotFound)
			return
		}

		isBlock := models.IsBlocked(ip)
		if isBlock {
			dc.OnCustomError(consts.ErrIpHasBlocked)
			return
		}
	}

	vid, err := extractors.MatchUrl(di.URL)
	if err != nil {
		dc.OnError(err)
		return
	}

	//dc.Ctx.ResponseWriter.Header().Set("Content-Type", "application/octet-stream")
	//dc.Ctx.ResponseWriter.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=a.mp4"))
	//dc.Ctx.ResponseWriter.Header().Set("Content-Length", fmt.Sprintf("%d", 1024000))

	file, _ := os.Create("temp.flv")
	header := map[string]string{}
	if len(vid.Formats) > 0 {
		for _, v := range vid.Formats[0].URLs {
			header["Referer"] = v.URL
			resp, err := utils.HttpGet(v.URL, header)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(err)
			writer := io.MultiWriter(file)
			io.Copy(writer, resp.Body)
		}
	}

	dc.JSON(vid)
}
