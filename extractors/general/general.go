package general

/*
	通用下载,用于尝试直接下载无法解析的 URL
*/

import (
	"strings"

	"git-me/common"
	"git-me/consts"
	"git-me/utils"
)

type BasicInfo struct {}

func (gn BasicInfo) Prepare(params map[string]interface{}) error {
	return nil
}

func (gn BasicInfo) Download(url string) (vid common.VideoData, err error) {
	exts := strings.Split(url, ".")
	ext := ""
	if len(exts) > 1 {
		ext = exts[len(exts)-1]
	}
	
	urlData := common.URLData{
		URL:  url,
		Size: consts.DefaultSize,
		Ext:  ext,
	}

	format := common.FormatData{URLs:[]common.URLData{urlData}}
	vid.Site = ""
	vid.Title = utils.FileName(exts[0])
	vid.Formats = []common.FormatData{format}

	return
}