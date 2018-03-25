package xiami

import (
	"encoding/xml"
	"fmt"

	"git-me/common"
	"git-me/utils"

	//"github.com/beevik/etree"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

type BasicInfo struct {
	Url    string
	Title  string
	Name   string
	LrcUrl string
}

type StringResources struct {
	XMLName        xml.Name         `xml:"resources"`
	ResourceString []ResourceString `xml:"string"`
}

type ResourceString struct {
	XMLName    xml.Name `xml:"string"`
	StringName string   `xml:"name,attr"`
	InnerText  string   `xml:",innerxml"`
}

func (xm BasicInfo) Prepare(Params map[string]interface{}) error {
	return nil
}

func (xm BasicInfo) Download(url string) (common.VideoData ,error) {
	// albums

	// collections

	// single track

	// mv
	return downloadMv(url)
}

func downloadSong(id, outputDir, infoOnly string) error {
	url := fmt.Sprintf("http://www.xiami.com/song/playlist/id/%s/object_name/default/object_id/0", id)
	_, err := utils.GetContent(url, common.FakeHeader)
	if err != nil {
		return err
	}
	return nil
}

func downloadMv(url string) (data common.VideoData,err error){
	page, err := utils.GetContent(url, nil)
	if err != nil {
		return
	}
	//fmt.Println(string(page))

	title := "xiami.flv"
	match := utils.Match(`<title>(^<]+)`, string(page))
	if len(match) > 0 {
		title = match[0]
	}

	data.Title = title
	data.Type = "video"

	vid, uid := "", ""
	match = utils.Match(`vid:"(\d+)"`, string(page))
	if len(match) > 0 {
		vid = strings.Split(match[0], `"`)[1]
		fmt.Println(vid)
	}

	match = utils.Match(`uid:"(\d+)"`, string(page))
	if len(match) > 0 {
		uid = strings.Split(match[0], `"`)[1]
		fmt.Println(uid)
	}

	apiUrl := fmt.Sprintf("http://cloud.video.taobao.com/videoapi/info.php?vid=%s&uid=%s", vid, uid)
	result, err := utils.GetContent(apiUrl, nil)
	if err != nil {
		return
	}

	fmt.Println(string(result))

	doc, err := goquery.NewDocument(apiUrl)
	if err != nil {
		return
	}

	str := doc.Find("video_url").Eq(-1).Text()
	end := doc.Find("length").Eq(-1).Text()
	str += fmt.Sprintf("/start_%d/end_%s/1.flv", 0, end)
	u := common.URLData{
		URL: str,
	}
	data.URLs = append(data.URLs, u)

	return
}
