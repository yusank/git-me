package youku

import (
	"fmt"
	"git-me/common"
	"git-me/utils"
	"log"
	"time"

	"github.com/bitly/go-simplejson"
)

type BasicInfo struct {
	Title string
	Url        string
	Referer    string
	Page       []byte
	VideoList  interface{}
	VideoNext  interface{}
	Password   string
	ApiData    *simplejson.Json
	ApiErrCode int
	ApiErrMsg  string

	Streams map[string]*StreamStruct

	CCode string
	Vid   string
	Utid  interface{}
	Ua    string

	PassProtected bool
}

type StreamStruct struct {
	Id           string
	Container    string
	VideoProfile string
	Size         int
	Pieces       []Piece
	M3u8Url      string
}

type Piece struct {
	Segs string
}

var StreamTypes map[string]*StreamStruct

func InitStream() {
	StreamTypes = map[string]*StreamStruct{
		"hd3":      &StreamStruct{Id: "hd3", Container: "flv", VideoProfile: "1080p"},
		"hd3v2":    &StreamStruct{Id: "hd3v2", Container: "flv", VideoProfile: "1080p"},
		"mp4hd3":   &StreamStruct{Id: "hd3v2", Container: "mp4", VideoProfile: "1080p"},
		"mp4hd3v2": &StreamStruct{Id: "hd3v2", Container: "mp4", VideoProfile: "1080p"},

		"hd2":   &StreamStruct{Id: "hd3", Container: "flv", VideoProfile: "超清"},
		"hd2v2": &StreamStruct{Id: "hd3v2", Container: "flv", VideoProfile: "超清"},

		// todo: 完善
	}
}

// something with cookies
//func (yk BasicInfo) FetchCna() {
//	var quotaCna = func(val string) string {
//		if strings.ContainsAny(val, "%") {
//			return val
//		}
//
//		return url.QueryEscape(val)
//	}
//}

func (yk BasicInfo) Ups() (err error) {
	url := fmt.Sprintf("https://ups.youku.com/ups/get.json?vid=%s&ccode=%d", yk.Vid, yk.CCode)
	url += "&client_ip=192.168.1.1"
	url += "&utid=" + fmt.Sprint(yk.Utid)
	url += "&client_ts" + fmt.Sprint(time.Now().Unix())
	if yk.PassProtected {
		url += "&password=" + yk.Password
	}
	headers := make(map[string]string)
	headers["Referer"] = yk.Referer
	headers["User-Agent"] = yk.Ua

	apiMate, err := utils.LoadJSON(url, headers)
	if err != nil {
		return
	}

	yk.ApiData = apiMate.Get("data")
	dataErr := yk.ApiData.Get("error")
	if dataErr != nil {
		yk.ApiErrCode, _ = dataErr.Get("code").Int()
		yk.ApiErrMsg, _ = dataErr.Get("note").String()
	}

	v := yk.ApiData.Get("videos")
	if v != nil {
		if l := v.Get("list"); l != nil {
			yk.VideoList = l
		}

		if n := v.Get("next"); n != nil {
			yk.VideoNext = n
		}
	}

	return
}

func (yk BasicInfo) GetVidFromUrl() error {
	b64p := `([a-zA-Z0-9=]+)`
	pList := []string{
		`youku\.com/v_shoe?id` + b64p,
		`player\.youku\.com/player\.php/sid/` + b64p + `\v/.swf`,
		`loader\.swf\?VideoIDS=` + b64p,
		`player\.youku\.com\.com/embed/` + b64p,
	}

	if yk.Url == "" {
		return common.ErrUrlIsEmpty
	}

	for _, v := range pList {
		results := utils.Match(v, yk.Url)
		if len(results) > 0 {
			yk.Vid = results[0]
			return nil
		}
	}

	return nil
}

func (yk BasicInfo) GetVidFromPage() (err error) {
	if yk.Url == "" {
		return common.ErrUrlIsEmpty
	}

	yk.Page, err = utils.GetContent(yk.Url, nil)
	if err != nil {
		return err
	}

	hit := utils.Match(`videoId2:"([A-Za-z0-9=]+)"`, string(yk.Page))
	if len(hit) >= 0 {
		yk.Vid = hit[0]
		return
	}

	return
}

func DownLoadByURL(url, outputDir string) {}

func (yk BasicInfo) Prepare(params map[string]interface{}) error {
	if yk.Url == "" && yk.Vid == "" {
		return nil
	}

	if yk.Url != "" && yk.Vid == "" {
		if err := yk.GetVidFromUrl(); err != nil {
			log.Println(err)
			return err
		}

		if yk.Vid == "" {
			if err := yk.GetVidFromPage(); err != nil {
				log.Println(err)
				return err
			}

			if yk.Vid == "" {
				log.Println("Cannot fatch vid")
				return nil
			}
		}
	}

	if a, found := params["src"]; found && a.(string) == "tudou" {
		yk.CCode = "0512"
	}

	if p, found := params["password"]; found {
		yk.PassProtected = true
		yk.Password = p.(string)
	}

	if err := yk.Ups(); err != nil {
		log.Println(err)
		return err
	}

	if yk.ApiData.Get("stream") == nil {
		if yk.ApiErrCode == -6001 { // wrong vid  parsed form page
			vidFromUrl := yk.Vid
			yk.GetVidFromPage()
			if vidFromUrl == yk.Vid {
				log.Println(yk.ApiErrMsg)
				return fmt.Errorf(yk.ApiErrMsg)
			}

			yk.Ups()
		}
	}

	if yk.ApiData.Get("stream") == nil {
		if yk.ApiErrCode == -2002 { // wrong password
			yk.PassProtected = true
			// it can be true already .offer another chance to retry
			yk.Password = utils.ReadInput("Password:")
			yk.Ups()
		}
	}

	if yk.ApiData.Get("stream") == nil {
		if yk.ApiErrMsg != "" {
			log.Fatal(yk.ApiErrMsg)
		} else {
			log.Fatal("unknown error")
		}
	}

	yk.Title, _ = yk.ApiData.Get("video").Get("title").String()

	streams, _ := yk.ApiData.Get("stream").Map()
	for _, v := range streams {
		streamId := v.(map[string]interface{})["stream_type"].(string)
		if a, found := StreamTypes[streamId]; found {
			streamId = a.Id

			if _,found := yk.Streams[streamId]; !found {
				yk.Streams[streamId] = &StreamStruct{
					Id:streamId,
					Container: a.Container,
					VideoProfile:a.VideoProfile,
					Size:streams["size"].(int),
					M3u8Url: streams["m3u8_url"].(string),
					Pieces: []Piece{
						Piece{Segs: streams["segs"].(string)},
					},
				}
			}
		}
	}

	return nil
}

func (yk BasicInfo) Download(params map[string]interface{}) error {
	return nil
}
