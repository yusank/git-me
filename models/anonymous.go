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

package models

import (
	"fmt"

	"github.com/yusank/git-me/db"
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

func IsBlocked(ip string) bool {
	has, err := db.Redis.SisMember(TimesRunOutUser, ip)
	if err != nil {
		return false
	}

	return has
}

func AddToDisabledUser(ip string) error {
	return db.Redis.SAdd(TimesRunOutUser, ip)
}
