package ipinfo

import (
	"goblin/pkg/ipinfo/geoip"
	"goblin/pkg/ipinfo/qqwry"
	"strings"
	log "unknwon.dev/clog/v2"
)

var DB Database

func Area(ip string) string {
	if DB != nil {
		return DB.Area(ip)
	}
	return ""
}

func (config *Config) Init() {
	switch config.Type {
	case "qqwry":
		DB = qqwry.New()
	case "geoip":
		DB = geoip.New(config.GeoLicenseKey)
	default:
		log.Fatal("wrong ipctl location database type: %q, type must %s", config.Type, strings.Join(ipType, ","))
	}
}
