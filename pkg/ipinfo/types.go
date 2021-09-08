package ipinfo

import (
	"fmt"
	"goblin/pkg/utils"
	"strings"
)

// 支持数据库类型
var ipType = []string{"qqwry", "geoip"}

type Config struct {
	Type          string `yaml:"type"`
	GeoLicenseKey string `yaml:"geo_license_key"`
}

// ValidateType 验证数据库配置信息
func (db *Config) ValidateType() error {
	if !utils.StrEqualOrInList(db.Type, ipType) {
		return fmt.Errorf("Ip db value %s type must %s ", db.Type, strings.Join(ipType, ","))
	}
	//
	if db.Type == "geoip" {
		if db.GeoLicenseKey == "" {
			return fmt.Errorf("geoip license_key is null")
		}
	}
	return nil
}
