package common

import (
	"fmt"
	"log"
)

type VideoExtractor interface {
	DownloadByUrl(vc *VideoCommon, url string, kv *map[string]interface{})
}

type Provider struct{}

func (vp *Provider) DownloadByUrl(vc *VideoCommon, url string, kv map[string]interface{}) {
	vc.Url = url
	if vc.Out {
		return
	}

}

func (vp *Provider) Download(vc *VideoCommon, kv map[string]interface{}) {
	if v, found := kv["infoOnly"]; found && v != nil {
		if v, found = kv["streamId"]; found && v != nil {
			if _, f := kv["index"]; f {
				vp.PrintI(v)
			} else {
				vp.Print(v)
			}
		}
	} else {
		streamId := ""
		if v, found = kv["streamId"]; found && v != nil {
			streamId = v.(string)
		} else {
			// chose best quality
			streamId = vc.StreamsSort[0].(string)
		}

		if v, found = kv["index"]; found {
			vp.PrintI(streamId)
		} else {
			vp.Print(streamId)
		}

		urls := []string{}
		ext := ""
		totalSize := 0
		if s, f := vc.Streams[streamId]; f {
			ss := s.(map[string]interface{})
			if u, f := ss["src"]; f {
				urls = u.([]string)
			}
			if e, f := ss["container"]; f {
				ext = e.(string)
			}
			if t, f := ss["size"]; f {
				totalSize = t.(int)
			}
		} else if ds, f := vc.DashStreams[streamId]; f {
			dss := ds.(map[string]interface{})
			if u, f := dss["src"]; f {
				urls = u.([]string)
			}
			if e, f := dss["container"]; f {
				ext = e.(string)
			}
			if t, f := dss["size"]; f {
				totalSize = t.(int)
			}
		}
		if ext == "m3u8" {
			ext = "mp4"
		}
		if urls == nil {
			log.Fatal("[Faild] Cannot extract video source")
		}

		header := make(map[string]string)
		if vc.UA != "" {
			header["User-Agent"] = vc.UA
		}
		if vc.Referer != "" {
			header["Referer"] = vc.Referer
		}
		DownloadURL(urls, []string{vc.Title}, ext, kv["outpuDir"].(string), totalSize,false, header)
	}
}

func (vp *Provider) Print(streamId interface{}) {
	fmt.Println(streamId)
}

func (vp *Provider) PrintI(streamId interface{}) {
	fmt.Printf("%v \n", streamId)
}
