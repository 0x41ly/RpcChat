package commons

import (
	"regexp"
)

func CheckUserName(username string) string {
	username = regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(username, "")
	return username
}
