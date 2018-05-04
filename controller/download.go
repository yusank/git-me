package controller

import (
	"encoding/json"

	"git-me/consts"
	"git-me/extractors"
	"git-me/models"
)

type DownloaderController struct {
	BasicController
}

type DownloadInfo struct {
	URL string `json:"url"`
}

func (dc *DownloaderController) ParseVideo() {
	var di DownloadInfo
	if err := json.Unmarshal(dc.Ctx.Input.RequestBody, &di); err != nil {
		dc.OnError(err)
		return
	}

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

	dc.JSON(vid)
}
