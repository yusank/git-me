package models

import (
	"fmt"
	"git-me/db"
)

const (
	AnonymousKey    = "anonymousUser:%s"
	TimesRunOutUser = "userTimesRunOut"

	KeyExpire = 60 * 60 * 24
	UseTimes  = 10
)

func AddToAnonUser(ip string) error {
	key := fmt.Sprintf(AnonymousKey, ip)
	count, err := db.Redis.Increase(key)
	if err != nil {
		return err
	}

	if count == UseTimes {
		if err := AddToDisabledUser(ip); err != nil {
			return err
		}

		return db.Redis.Delete(key)
	}

	return nil
}

func IsDisabled(ip string) bool {
	has, err := db.Redis.SisMember(TimesRunOutUser, ip)
	if err != nil {
		return false
	}

	return has
}

func AddToDisabledUser(ip string) error {
	return db.Redis.SAdd(TimesRunOutUser, ip)
}
