package models

import "git-me/db"

const (
	StatSitePrefix = "StatParseNum"
)

func AddStatSite(site string) error {
	return db.Redis.ZIncrby(StatSitePrefix, 1, site)
}
