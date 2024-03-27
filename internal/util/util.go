package util

import (
	"net/url"
)

// ValidateAddress проверяем валидность URL
func ValidateAddress(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
