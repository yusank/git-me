package utils

import "regexp"

// IsValidAccount - 账户名是否合法
func IsValidAccount(name string) bool {
	match, _ := regexp.MatchString("^[\\w]{6,127}$", name)
	return match
}

// IsValidPassword - 密码是否合法
func IsValidPassword(pass string) bool {
	match, _ := regexp.MatchString("^[\\w#&]{6,127}$", pass)
	return match
}

const (
	mailPattern = `^[a-z0-9A-Z]+([\-_\.][a-z0-9A-Z]+)*@([a-z0-9A-Z]+(-[a-z0-9A-Z]+)*?\.)+[a-zA-Z]{2,4}$`
)

func IsValidEmail(email string) bool {
	if len(email) < 6 {
		return false
	}

	return regexp.MustCompile(mailPattern).MatchString(email)
}
