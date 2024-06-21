package app

import (
	"fmt"
	"github.com/vadskev/urlshort/config"
	"github.com/vadskev/urlshort/internal/util"
)

var URLDBase = make(map[string]string)

// ShortURL
func ShortURL(inputURL string) (string, error) {
	if len(inputURL) < 1 {
		return "", fmt.Errorf("the URL must not be empty")
	}

	shortCode := util.GenerateRandomString()
	URLDBase[shortCode] = inputURL

	return config.GetConfig().HostServer + "/" + shortCode, nil
}

// ExpandURL
func ExpandURL(code string) (string, error) {
	if len([]rune(code)) != 7 {
		return "", fmt.Errorf("invalid url id")
	}

	expandedURL, ok := URLDBase[code]
	if !ok {
		return "", fmt.Errorf("url does not exist")
	}

	return expandedURL, nil
}
