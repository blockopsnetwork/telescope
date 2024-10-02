package networks

import (
	"net/url"
	"strings"
)

type ScrapeConfig struct {
	JobName       string
	StaticConfigs []StaticConfig
}

type StaticConfig struct {
	Targets []string
}

func toLowerAndEscape(input string) string {
	lowercase := strings.ToLower(input)
	escaped := url.QueryEscape(lowercase)
	return escaped
}
