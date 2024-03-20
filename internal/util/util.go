package util

import (
	"net/url"
)

// IsURL проверяем валидность URL
func IsURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
