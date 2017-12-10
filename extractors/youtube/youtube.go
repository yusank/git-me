package youtube

import "git-me/common"

type Youtube struct {
	Name    string
	Handler *common.Provider
}

func (ys *Youtube) Prepare(info *common.VideoCommon, kv map[string]interface{}) {
	// do something
}
