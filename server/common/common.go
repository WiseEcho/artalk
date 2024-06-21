package common

import (
	"strings"

	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/ArtalkJS/Artalk/internal/utils"
)

const (
	UrlTimeout = 1800
)

func GetOpenURLByKey(urlKey string, cfg *config.Config) string {
	if urlKey == "" {
		return ""
	}
	openURL := cfg.FS + urlKey
	if IsUploadUrlKey(urlKey) {
		sign, _ := utils.EncryptURL(urlKey, UrlTimeout)
		openURL += "?sign=" + sign
	}
	return openURL
}

func IsUploadUrlKey(urlKey string) bool {
	if urlKey == "" {
		return false
	}
	if strings.HasPrefix(urlKey, "/web/") || strings.HasPrefix(urlKey, "/mini_web/") {
		return true
	}
	return false
}
